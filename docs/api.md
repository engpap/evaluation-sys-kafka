# General Requirements
- 3 type of users: students, professors, admins
- 4 micro-services: users, courses, projects, registration
- No database; implement state using in-memory data structures.

# Models

## Student
- ID

## Professor
- ID

## Project
- ID
- CourseID
- Name

## Grade
- SubmissionID
- ProfessorID
- Grade

## Course
- ID
- Name

## Enrollment
- StudentID
- CourseID

## Submission
- ID
- StudentID
- ProjectID

## Completed Course
- ID
- StudendID
- CourseID

# Endpoints

## User Service
- POST: /users/stud/create -> just REST api ✅
    - student id   
- POST: /users/prof/create -> just REST api ✅
    - professor id
- GET: /users/stud/carreer -> REST api + fetch from in-state memory populated through consumer on `completed-course` topic 

**Additional:**
- GET: /users/stud/:stud-id/enrollments -> REST api + fetch from in-state memory populated through consumer on `enrollment` topic
- GET: /users/prof/:prof-id/submissions -> REST api + fetch from in-state memory populated through consumer on `submission` topic


## Course Service
- POST: /courses/create -> just REST api ✅
    - course id
    - name
- DELETE: /courses/:course-id/delete -> just REST api ✅
- POST: /courses/enroll -> REST api (+ send on `enrollment` topic) + fetch from in-state memory populated through consumer on `student` topic to check that student id exist ✅
    - student id
    - course id

**Additional:**
- GET: /courses/:course-id -> REST api + fetch from in-state memory populated through consumer on `project` topic


## Project Service
- POST: /projects/create -> REST api + send on `project` topic ✅
    - project id
    - course id
    - name
- POST: /projects/:project-id/submit -> REST api (+ send on `submission` topic)✅
    - submission id
    - student id
    - project id
    TODO: here you can listen on `enrollment` and check whether the student that's is submitting is actually attending the course
- POST: /projects/:project-id/submission-id={id}/grade -> REST api + send on `grade` topic ✅
    - submission id
    - professor id
    - grade


## Registration Service
No endpoints.
Consumes on `grade`, `course`, `project` topic.
Calculates whether a course is completed for a student, i.e. if the student delivered all projects for that course and the sum of the grades is suffcient.
Produces to `completed-course` topic.

# Assumptions
1. Login of users is not implemented

# Legend/Key
Additional features are not manadatory. They could be implemented to make the app more complete.
Phrases enclosed in `()` are not mandatory, they're used for implementing additional features.

# TODO
- when you create sub, specify the course id, not only project id
- when sub  is created, check if student exists
- check if student enrolled in course before submitting a proj
- finish frontend cli
- implement faulty recovery procedure


# Fault Recovery
Implement a fault recovery procedure to resume a valid state of the services.

## Scenario 1
other services are running, one goes down, a producer sends messages on the topic on which the failed service was listening, then the failed service is relaunched.
The consumer correctly re-consumes the messages from the beginning of the history (seek to the beginning of the assigned partition) => State is recovered. ✅
This is done thanks to `resetPartitionsOffset` function.

## Scenario 2
3/4 services are running, they produce messages on the topic on which service 1/4 should be consuming from. Service 1/4 is launched. It correctly read the messages that have been produced. ✅
This is done automatically by Kafka with the current config of the consumer.