FROM alpine
ADD main /
ADD .env /
#EXPOSE 8080
EXPOSE 80

CMD ["/main"]