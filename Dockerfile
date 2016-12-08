FROM scratch
ADD main /
ADD .env /
EXPOSE 8080

CMD ["/main"]