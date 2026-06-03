pipeline {
    agent any

    environment {
        APP_SRC_DIR    = "app_src"
        COMPOSE_FILE   = "docker-compose-app.yaml"
        IMAGE_TAG      = "latest"
        TRIVY_OUTPUT_DIR = "trivy-reports"
        SonarQube_Project = "NT548-Lab-Nhom10"
    }

    stages {

        // ─────────────────────────────────────────────
        // STAGE 1: Checkout
        // ─────────────────────────────────────────────
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        // ─────────────────────────────────────────────
        // STAGE 2: SonarQube Scan
        // ─────────────────────────────────────────────
        stage('SonarQube Scan') {
            steps {
                withSonarQubeEnv('sonarqube-server') {
                    script {
                        def scannerHome = tool 'sonar-scanner'
                        sh """
                            ${scannerHome}/bin/sonar-scanner \
                                -Dsonar.projectKey=${SonarQube_Project} \
                                -Dsonar.sources=${APP_SRC_DIR} \
                                -Dsonar.java.binaries=. \
                                -Dsonar.exclusions=**/*.java
                        """
            }
                }
            }
            post {
                failure {
                    echo 'SonarQube scan failed, but pipeline continues.'
                }
            }
        }

        // ─────────────────────────────────────────────
        // STAGE 3: Trivy Scan - Dockerfiles
        // ─────────────────────────────────────────────
        stage('Trivy Scan Dockerfiles') {
            steps {
                script {
                    sh "mkdir -p ${TRIVY_OUTPUT_DIR}"

                    def services = sh(
                        script: "ls -d ${APP_SRC_DIR}/*/",
                        returnStdout: true
                    ).trim().split('\n')

                    for (svc in services) {
                        def svcName    = svc.replaceAll('/$', '').split('/').last()
                        def dockerfile = "${svc}Dockerfile"
                        def reportFile = "${TRIVY_OUTPUT_DIR}/dockerfile-${svcName}.txt"

                        if (fileExists(dockerfile)) {
                            echo "Scanning Dockerfile: ${dockerfile}"
                            sh """
                                trivy config \
                                    --exit-code 0 \
                                    --output ${reportFile} \
                                    ${svc}
                            """
                        } else {
                            echo "No Dockerfile found in ${svc}, skipping."
                        }
                    }
                }
            }
        }

        // ─────────────────────────────────────────────
        // STAGE 4: Docker Build
        // ─────────────────────────────────────────────
        stage('Docker Build') {
            steps {
                script {
                    def services = sh(
                        script: "ls -d ${APP_SRC_DIR}/*/",
                        returnStdout: true
                    ).trim().split('\n')

                    for (svc in services) {
                        def svcName    = svc.replaceAll('/$', '').split('/').last()
                        def dockerfile = "${svc}Dockerfile"

                        if (fileExists(dockerfile)) {
                            echo "Building image: ${svcName}:${IMAGE_TAG}"
                            sh """
                                docker build \
                                    -t ${svcName}:${IMAGE_TAG} \
                                    -f ${dockerfile} \
                                    ${svc}
                            """
                        } else {
                            echo "No Dockerfile found in ${svc}, skipping build."
                        }
                    }
                }
            }
        }

        // ─────────────────────────────────────────────
        // STAGE 5: Trivy Scan - Docker Images
        // ─────────────────────────────────────────────
        stage('Trivy Scan Images') {
            steps {
                script {
                    def services = sh(
                        script: "ls -d ${APP_SRC_DIR}/*/",
                        returnStdout: true
                    ).trim().split('\n')

                    for (svc in services) {
                        def svcName    = svc.replaceAll('/$', '').split('/').last()
                        def dockerfile = "${svc}Dockerfile"
                        def reportFile = "${TRIVY_OUTPUT_DIR}/image-${svcName}.txt"

                        if (fileExists(dockerfile)) {
                            echo "Scanning image: ${svcName}:${IMAGE_TAG}"
                            sh """
                                trivy image \
                                    --exit-code 0 \
                                    --output ${reportFile} \
                                    ${svcName}:${IMAGE_TAG}
                            """
                        }
                    }
                }
            }
        }

        // ─────────────────────────────────────────────
        // STAGE 6: Deploy with Docker Compose
        // ─────────────────────────────────────────────
        stage('Deploy Docker Compose') {
            steps {
                sh """
                    docker compose -f ${COMPOSE_FILE} down --remove-orphans || true
                    docker compose -f ${COMPOSE_FILE} up -d --no-build
                """
            }
        }

    }

    // ─────────────────────────────────────────────
    // POST: Archive Trivy reports
    // ─────────────────────────────────────────────
    post {
        always {
            archiveArtifacts artifacts: "${TRIVY_OUTPUT_DIR}/**", allowEmptyArchive: true
            echo "Pipeline finished. Trivy reports archived."
        }
        success {
            echo "All stages completed successfully."
        }
        failure {
            echo "Pipeline failed. Check logs above."
        }
    }
}