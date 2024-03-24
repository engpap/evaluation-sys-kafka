# Project #5: Online services for continuous evaluation

## Description
In this project, you are to implement the online services that support university courses adopting continuous evaluation. The application consists of a frontend that accepts requests from users and a backend that processes them.

There are three types of users interacting with the service:
1. students enroll for courses, submit their solutions for projects, and check the status of their submissions;
2. admins can add new courses and remove old courses;
3. professors post projects and grade the solutions submitted by students.

The backend is decomposed in four services, following the microservices paradigm:
1. the users service manages personal data of registered students and professors;
2. the courses service manages courses, which consists of several projects;
3. the projects service manages submission and grading of individual projects;
4. the registration service handles registration of grades for completed courses. A course is completed for a student if the student delivered all projects for that course and the sum of the grades is sufficient.

## Assumptions and Guidelines
1. Services do not share state, but only communicate by exchanging messages/events over Kafka topics
    ○ They adopt an event-driven architecture: you can read chapter 5 of the book “Designing Event-Driven Systems”[1]to get some design principles for your project
2. Services can crash at any time and lose their state
    ○ You may simply implement state using in-memory data structures
    ○ You need to implement a fault recovery procedure to resume a valid state of the services
3. You can assume that Kafka topics cannot be lost
4. You can use any technology to implement the frontend (e.g., simple command-line application or basic
REST/Web app)