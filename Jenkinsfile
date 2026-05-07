pipeline {
    agent any
    
    environment {
        DOCKER_REGISTRY = 'your-registry.com'  // 修改为你的镜像仓库地址
        IMAGE_TAG = "${BUILD_NUMBER}"
        BACKEND_IMAGE = "intehub-backend:${IMAGE_TAG}"
        FRONTEND_IMAGE = "intehub-frontend:${IMAGE_TAG}"
        
        // 从 Jenkins Credentials 获取敏感信息
        POSTGRESQL_URI = credentials('intehub-postgresql-uri')  // 数据库连接URI
        JWT_SECRET = credentials('intehub-jwt-secret')          // JWT 密钥
        CRYPTO_KEY = credentials('intehub-crypto-key')         // 加密密钥
    }
    
    stages {
        stage('Checkout') {
            steps {
                echo '拉取代码...'
                checkout scm
            }
        }
        
        stage('Build Backend') {
            steps {
                echo '构建后端镜像...'
                script {
                    sh '''
                        docker build -t ${BACKEND_IMAGE} -f Dockerfile.backend .
                    '''
                }
            }
        }
        
        stage('Build Frontend') {
            steps {
                echo '构建前端镜像...'
                script {
                    sh '''
                        docker build -t ${FRONTEND_IMAGE} -f Dockerfile.frontend .
                    '''
                }
            }
        }
        
        stage('Deploy') {
            steps {
                echo '部署应用...'
                script {
                    // 停止并删除旧容器
                    sh '''
                        docker stop intehub-backend intehub-frontend || true
                        docker rm intehub-backend intehub-frontend || true
                    '''
                    
                    // 创建网络（如果不存在）
                    sh 'docker network create intehub-network || true'
                    
                    // 启动后端容器
                    sh '''
                        docker run -d \
                          --name intehub-backend \
                          --network intehub-network \
                          -e INTEHUB_POSTGRESQL_URI="${POSTGRESQL_URI}" \
                          -e INTEHUB_SERVER_PORT=8080 \
                          -e INTEHUB_JWT_SECRET="${JWT_SECRET}" \
                          -e INTEHUB_CRYPTO_KEY="${CRYPTO_KEY}" \
                          -p 8080:8080 \
                          --restart unless-stopped \
                          ${BACKEND_IMAGE}
                    '''
                    
                    // 等待后端启动
                    sh 'sleep 5'
                    
                    // 启动前端容器
                    sh '''
                        docker run -d \
                          --name intehub-frontend \
                          --network intehub-network \
                          -p 80:80 \
                          --restart unless-stopped \
                          ${FRONTEND_IMAGE}
                    '''
                }
            }
        }
        
        stage('Health Check') {
            steps {
                echo '检查服务健康状态...'
                script {
                    // 等待服务启动
                    sh 'sleep 10'
                    
                    // 检查后端
                    def backendHealth = sh(
                        script: "curl -f -s http://localhost:8080/api/health || exit 1",
                        returnStatus: true
                    )
                    
                    if (backendHealth != 0) {
                        error "后端健康检查失败"
                    }
                    
                    // 检查前端
                    def frontendHealth = sh(
                        script: "curl -f -s http://localhost:80 || exit 1",
                        returnStatus: true
                    )
                    
                    if (frontendHealth != 0) {
                        error "前端健康检查失败"
                    }
                    
                    echo '所有服务健康检查通过！'
                }
            }
        }
        
        stage('Clean Old Images') {
            steps {
                echo '清理旧镜像...'
                script {
                    sh '''
                        # 清理后端旧镜像（保留最新3个版本）
                        docker images intehub-backend --format "{{.Tag}}" | \
                        grep -E '^[0-9]+$' | \
                        sort -rn | \
                        tail -n +4 | \
                        xargs -r -I {} docker rmi intehub-backend:{} || true
                        
                        # 清理前端旧镜像（保留最新3个版本）
                        docker images intehub-frontend --format "{{.Tag}}" | \
                        grep -E '^[0-9]+$' | \
                        sort -rn | \
                        tail -n +4 | \
                        xargs -r -I {} docker rmi intehub-frontend:{} || true
                        
                        # 清理悬空镜像
                        docker image prune -f
                    '''
                }
            }
        }
    }
    
    post {
        success {
            echo '✅ 部署成功！'
        }
        failure {
            echo '❌ 部署失败！开始回滚...'
            script {
                sh '''
                    docker stop intehub-backend intehub-frontend || true
                    docker rm intehub-backend intehub-frontend || true
                '''
            }
        }
        always {
            echo '清理工作空间...'
            cleanWs()
        }
    }
}
