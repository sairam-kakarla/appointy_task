# Appointy Task
## Description
You are required to Design and Develop an HTTP JSON API capable of the following operations,
- Create an User
  - Should be a POST request
  - Use JSON request body
  - URL should be ‘/users'
- Get a user using id
  - Should be a GET request
  - Id should be in the url parameter
  - URL should be ‘/users/<id here>’
- Create a Post
  - Should be a POST request
  - Use JSON request body
  - URL should be ‘/posts'
- Get a post using id
  - Should be a GET request
  - Id should be in the url parameter
  - URL should be ‘/posts/<id here>’
- List all posts of a user
  - Should be a GET request
  - URL should be ‘/posts/users/<Id here>'

## Modules required
  - [MongoDB Go Driver](https://github.com/mongodb/mongo-go-driver)  
  ```
  go get go.mongodb.org/mongo-driver/mongo
  ```
## Installation
Download and install [Go](https://golang.org/)
```
  git clone https://github.com/sairam-kakarla/appointy_task
  cd appointy_task
  go mod init
  go build main.go
  ./main.go
```
## Usage
- Endpoint ```users```
  - GET  ```/users/<id>``` to retrieve user information
  - POST ```/users``` with http request body containing the user information
- Endpoint ```/posts```
  - GET ```posts/<id>``` to retrieve post information
  -POST ```/post``` with http request body containing the post information
- Endpoint ```/posts/users```
  - GET ```/posts/users/<id>``` to retrieve posts posted by user 'id'.
  
