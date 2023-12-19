FROM amd64/alpine
EXPOSE 3001
WORKDIR /app
COPY account_backend .
CMD ["/app/account_backend"]
