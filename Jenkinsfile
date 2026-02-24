pipeline {
    agent any
    
    environment {
        DOCKER_REGISTRY = 'your-registry.com'  // 修改为你的镜像仓库地址
        IMAGE_TAG = "${BUILD_NUMBER}"
        BACKEND_IMAGE = "intehub-backend:${IMAGE_TAG}"
        FRONTEND_IMAGE = "intehub-frontend:${IMAGE_TAG}"
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
                        docker tag ${BACKEND_IMAGE} ${BACKEND_IMAGE}
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
                        docker tag ${FRONTEND_IMAGE} ${FRONTEND_IMAGE}
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
                        docker stop intehub-backend intehub-frontend intehub-postgres || true
                        docker rm intehub-backend intehub-frontend intehub-postgres || true
                    '''
                    
                    // 创建网络（如果不存在）
                    sh 'docker network create intehub-network || true'
                    
                    // 启动 PostgreSQL 容器
                    sh '''
                        docker run -d \
                          --name intehub-postgres \
                          --network intehub-network \
                          -e POSTGRES_DB=intehub \
                          -e POSTGRES_USER=intehub \
                          -e POSTGRES_PASSWORD=intehub123 \
                          -v intehub-postgres-data:/var/lib/postgresql/data \
                          --restart unless-stopped \
                          postgres:16-alpine
                    '''
                    
                    // 等待数据库启动
                    sh 'sleep 10'
                    
                    // 启动后端容器
                    sh '''
                        docker run -d \
                          --name intehub-backend \
                          --network intehub-network \
                          -e INTEHUB_POSTGRESQL_URI="host=intehub-postgres port=5432 user=intehub password=intehub123 dbname=intehub sslmode=disable" \
                          -e INTEHUB_SERVER_PORT=8080 \
                          -e INTEHUB_JWT_SECRET="your-secret-key-change-in-production" \
                          -p 8080:8080 \
                          --restart unless-stopped \
                          ${BACKEND_IMAGE}
                    '''
                    
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
            echo '部署成功！'
        }
        failure {
            echo '部署失败！'
        }
        always {
            echo '清理工作空间...'
            cleanWs()
        }
    }
}
