from alpine:3.6

ARG TOKEN 

ENV TOKEN=$TOKEN

COPY beermate /

RUN chmod 755 /beermate

RUN apk add ca-certificates --update

CMD "./beermate" $TOKEN
