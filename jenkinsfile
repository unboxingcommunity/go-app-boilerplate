pipeline {
    agent any

    stages {
        stage('Build Docker Image') {
            steps {
                script {
                    git 'https://github.com/unboxingcommunity/go-app-boilerplate.git'

                    dir('docker') {
                        docker.build('dockerfile:latest')
                    }
                }
            }
        }
    }
}
