# Finbolt-Lab User Service Built with GO

## Getting Started.

Follow these instructions to set up and run the web server on your local machine

### Prerequisites

[Go](https://golang.org/) installed on your machine

### Installation

1. Clone the repository:
```sh
    git clone https://github.com/Ebentim/finbolt-user-service.git
```

2. Download the dependencies
```go
    go mod tidy
```
### Running the Server

1. Run the web server:
```go
    go run main.go
```
Open your web browser and navigate to http://localhost:5080 to see the web server in action.

### TODO

1. Security
    
        - Verify user tokens with firebase admin
        - Encrypt user data

2. Storage

        - Change image storage from base64 to a more efficient less memory intensive format
        - Add other data table as needed e.g, personalized currriculum forcasting, simulation results etc
        -role based learning
        -Categorize transactions intelligently using ai (Obtain from Learning Service)

        

3. Business

        - Ai for daily flash cards
        - Cashing
        - Scaleling