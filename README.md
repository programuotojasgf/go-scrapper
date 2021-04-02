<p align="center">
<img width="250" src="https://i.ibb.co/RbkqK08/scrapper-with-plastic-handle-500x500.jpg" />
</p>

![alt text](http://url/to/img.png)

<p align="center">
  <a href="https://goreportcard.com/report/github.com/programuotojasgf/go-scrapper">
    <img src="https://goreportcard.com/badge/github.com/programuotojasgf/go-scrapper">
  </a>
</p>
  

# Shopify review scrapper

go-scrapper is a scrapping for shopify reviews. The output is a database where all the 3 word phrase frequency from reviews is aggregated. A REST API to retrieve the data is a standalone project which can be found at [https://github.com/programuotojasgf/go-server](https://github.com/programuotojasgf/go-server)

## Installation

1. Install go. An installation guide can be found at [https://golang.org/doc/install](https://golang.org/doc/install)
2. Setup the configuration file in `config/config.development.json`  

Notes: 
* The application will create the database and collection by itself, if they do not already exist. 
* You will have to provide a MongoDB server connection.
* The was created using GoLand 2020.3.4, but it's not required to use the infrastructure for it.

## Usage

Open the terminal window, navigate to the project directory and execute the following
```console
go run .
```
This will launch the application. It will start scrapping and pushing the results to the database. The process is loged, so you can see what's happening. Once the application finished scrapping all the reviews it will print out `Finished processing all review contents.` and exit.

You can launch the application multiple times, to update the data or to resume at a later point - consumed reviews will not be re-processed, only new ones.

## Troubleshooting

If you encounter issues with GOROOT using GoLand, go to [https://www.jetbrains.com/help/go/configuring-goroot-and-gopath.html](https://www.jetbrains.com/help/go/configuring-goroot-and-gopath.html) for tips how to solve this issue.