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
                    sh '''
                        export IMAGE_TAG=${IMAGE_TAG}
                        docker-compose down
                        docker-compose up -d
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
