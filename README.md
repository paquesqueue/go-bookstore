# Description

    Bookstore เป็นโปรเจ็คที่สร้างขึ้นเพื่อพัฒนาระบบ Backend Web Service อย่างง่าย โดยมีระบบ user สำหรับการจัดการข้อมูลผู้ใช้งานและระบบ book สำหรับการจัดการหนังสือ 
    
    ใช้ภาษา Go เป็นภาษาหลัก นำ Echo ที่เป็น Web Framework ของภาษา Go มาพัฒนาระบบ API เพื่อเชื่อมต่อกับฐานข้อมูล PostgreSQL นอกจากนี้มีการใช้ Docker ในการจำลองสภาพแวดล้อม เพื่อสร้าง container สำหรับ run ตัว app และตัว database และนำมาใช้เป็น sandbox สำหรับ Integration Test

# Technology Stack

    Programming Language : Go (Echo Web Framework)
    Database : PostgreSQL
    Tools: Git, Docker

# Technical Requirement

    Go 1.19

# Instructions

# Run Application Locally

        $ DRIVER_NAME=postgres DATABASE_URL=postgres://<database_url>?sslmode=disable PORT=<port> ACCESS_TOKEN=token go run main.go
        
# Run Application in Container on Docker

        ดู Makefile ประกอบ

            $ docker build -t go-bookstore:latest .
         
            $ docker compose -f docker-compose-postgres.yml up --detach

            $ docker compose -f docker-compose.yml up --detach

            $ DRIVER_NAME=postgres DATABASE_URL=postgres://user:p@ssw0rd@localhost:5432/go-bookstore-db?sslmode=disable PORT=2565 ACCESS_TOKEN=token go run main.go

        * Stop the Running Containers
         
            $ docker compose -f docker-compose-postgres.yml stop

            $ docker compose -f docker-compose.yml stop

# Run Unit Testing

            $ go clean -testcache && go test -v --tags=unit ./...

# Run Integration Test on Docker
   
            $ docker compose -f docker-compose.test.yml up --build --abort-on-container-exit --exit-code-from it_tests