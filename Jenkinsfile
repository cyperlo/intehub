pipeline {
    agent any
    
    environment {
        // 镜像标签（本地构建，不需要推送）
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
        
        stage('Build Images') {
            parallel {
                stage('Build Backend') {
                    steps {
                        echo '构建后端镜像...'
                        script {
                            sh '''
                                docker build -t ${BACKEND_IMAGE} -f Dockerfile.backend .
                                docker tag ${BACKEND_IMAGE} intehub-backend:latest
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
                                docker tag ${FRONTEND_IMAGE} intehub-frontend:latest
                            '''
                        }
                    }
                }
            }
        }
        
        stage('Deploy') {
            steps {
                echo '部署应用...'
                script {
                    // 创建网络（如果不存在）
                    sh 'docker network create intehub-network || true'
                    
                    // 停止旧容器（优雅停止，给30秒时间）
                    sh '''
                        if docker ps -a | grep -q intehub-backend; then
                            echo "停止旧的后端容器..."
                            docker stop -t 30 intehub-backend || true
                            docker rm intehub-backend || true
                        fi
                        
                        if docker ps -a | grep -q intehub-frontend; then
                            echo "停止旧的前端容器..."
                            docker stop -t 30 intehub-frontend || true
                            docker rm intehub-frontend || true
                        fi
                    '''
                    
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
                          -v /var/intehub/data:/app/data \
                          -v /var/intehub/logs:/app/logs \
                          --restart unless-stopped \
                          --log-driver json-file \
                          --log-opt max-size=10m \
                          --log-opt max-file=3 \
                          ${BACKEND_IMAGE}
                    '''
                    
                    echo '等待后端服务启动...'
                    sh 'sleep 10'
                    
                    // 启动前端容器
                    sh '''
                        docker run -d \
                          --name intehub-frontend \
                          --network intehub-network \
                          -p 80:80 \
                          --restart unless-stopped \
                          --log-driver json-file \
                          --log-opt max-size=10m \
                          --log-opt max-file=3 \
                          ${FRONTEND_IMAGE}
                    '''
                    
                    echo '容器启动完成'
                }
            }
        }
        
        stage('Health Check') {
            steps {
                echo '检查服务健康状态...'
                script {
                    // 检查容器是否运行
                    def backendRunning = sh(
                        script: "docker ps | grep intehub-backend",
                        returnStatus: true
                    )
                    
                    if (backendRunning != 0) {
                        echo "后端容器日志："
                        sh 'docker logs intehub-backend || true'
                        error "后端容器未运行"
                    }
                    
                    def frontendRunning = sh(
                        script: "docker ps | grep intehub-frontend",
                        returnStatus: true
                    )
                    
                    if (frontendRunning != 0) {
                        echo "前端容器日志："
                        sh 'docker logs intehub-frontend || true'
                        error "前端容器未运行"
                    }
                    
                    // 等待服务完全启动
                    echo '等待服务完全启动...'
                    sh 'sleep 5'
                    
                    // 检查后端健康接口
                    retry(3) {
                        sh 'sleep 2'
                        def backendHealth = sh(
                            script: "curl -f -s http://localhost:8080/api/health || exit 1",
                            returnStatus: true
                        )
                        
                        if (backendHealth != 0) {
                            error "后端健康检查失败，重试中..."
                        }
                    }
                    
                    // 检查前端
                    def frontendHealth = sh(
                        script: "curl -f -s -o /dev/null -w '%{http_code}' http://localhost:80",
                        returnStdout: true
                    ).trim()
                    
                    if (frontendHealth != '200') {
                        error "前端健康检查失败，HTTP状态码: ${frontendHealth}"
                    }
                    
                    echo '✅ 所有服务健康检查通过！'
                }
            }
        }
        
        stage('Show Logs') {
            steps {
                echo '显示容器启动日志...'
                script {
                    sh '''
                        echo "=== 后端容器日志（最后20行）==="
                        docker logs --tail 20 intehub-backend
                        
                        echo ""
                        echo "=== 前端容器日志（最后20行）==="
                        docker logs --tail 20 intehub-frontend
                    '''
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
            echo "后端服务: http://localhost:8080"
            echo "前端服务: http://localhost:80"
        }
        failure {
            echo '❌ 部署失败！开始回滚...'
            script {
                sh '''
                    # 停止失败的容器
                    docker stop intehub-backend intehub-frontend || true
                    docker rm intehub-backend intehub-frontend || true
                    
                    # 尝试启动上一个版本（如果存在）
                    LAST_BUILD=$((BUILD_NUMBER - 1))
                    if docker images intehub-backend:${LAST_BUILD} | grep -q ${LAST_BUILD}; then
                        echo "回滚到版本 ${LAST_BUILD}"
                        docker run -d \
                          --name intehub-backend \
                          --network intehub-network \
                          -e INTEHUB_POSTGRESQL_URI="${POSTGRESQL_URI}" \
                          -e INTEHUB_SERVER_PORT=8080 \
                          -e INTEHUB_JWT_SECRET="${JWT_SECRET}" \
                          -e INTEHUB_CRYPTO_KEY="${CRYPTO_KEY}" \
                          -p 8080:8080 \
                          --restart unless-stopped \
                          intehub-backend:${LAST_BUILD}
                        
                        docker run -d \
                          --name intehub-frontend \
                          --network intehub-network \
                          -p 80:80 \
                          --restart unless-stopped \
                          intehub-frontend:${LAST_BUILD}
                    fi
                '''
            }
        }
        always {
            echo '清理工作空间...'
            cleanWs()
        }
    }
}
