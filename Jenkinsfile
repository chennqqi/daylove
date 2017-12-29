pipeline {
    agent any

    stages {
        stage('Build') {
            steps {
                echo 'Building..'
                sh "docker build -t daylove ."
            }
        }
        stage('Test') {
            steps {
                echo 'Testing..'
            }
        }
        stage('Deploy') {
            steps {
                echo 'Deploying....'
                sh "rsync -avzP ./* /data/www/daylove/"                
                sh "cd /data/www/web-svc && ls -al && ./force-replace.sh daylove &&./reload.sh"
            }
        }
    }
}

