FROM golang:alpine

WORKDIR /quics

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY . .
RUN go build -o qis ./cmd

ENV PATH="/quics:${PATH}"

EXPOSE 6121/udp
EXPOSE 6122/udp

CMD [ "qis", "run"]