pipeline {
    agent any

    stages {
        stage('Build Docker Image') {
            steps {
                checkout([$class: 'GitSCM', branches: [[name: '*/build-Branch']], doGenerateSubmoduleConfigurations: false, extensions: [], submoduleCfg: [], userRemoteConfigs: [[url: 'https://github.com/unboxingcommunity/go-app-boilerplate.git']]])

                script {
                    // Additional Docker build command
                    sh 'docker build -f docker/Dockerfile .'

                    // Continue with the original Docker build step
                    dir('docker') {
                        sh 'docker build -f Dockerfile .'
                    }
                }
            }
        }
    }
}

