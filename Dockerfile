FROM golang:1.8
WORKDIR /go/src/app
COPY . .

#Go
RUN go-wrapper download   # "go get -d -v ./..."
RUN go-wrapper install    # "go install -v ./..."

#Install Node
RUN apt-get update -y
RUN apt-get install -y node \
    npm
RUN npm install --prefix ./app -g

CMD ["go-wrapper", "run"] # ["app"]
CMD ["npm", "start"]
EXPOSE 8080