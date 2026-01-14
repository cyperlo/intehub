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
                        docker stop intehub-backend || true
                        docker rm intehub-backend || true
                        docker stop intehub-frontend || true
                        docker rm intehub-frontend || true
                    '''
                    
                    // 创建网络（如果不存在）
                    sh 'docker network create intehub-network || true'
                    
                    // 启动后端容器
                    sh '''
                        docker run -d \
                          --name intehub-backend \
                          --network host \
                          -v $(pwd)/config.yaml:/app/config.yaml \
                          -v $(pwd)/data:/app/data \
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
