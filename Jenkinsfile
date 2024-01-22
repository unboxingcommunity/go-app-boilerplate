pipeline {
    agent any

    stages {
        stage('Build Docker Image') {
            steps {
                script {
                    git 'https://github.com/unboxingcommunity/go-app-boilerplate.git'

                    dir('docker') {
                        sh 'docker build -f Dockerfile .'
                    }
                }
            }
        }
    }
}

