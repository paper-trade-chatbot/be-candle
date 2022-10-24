FROM alpine:latest

ENV SERVER_LISTEN_ADDRESS '0.0.0.0'
ENV SERVER_LISTEN_PORT '80'
ENV GRPC_SERVER_LISTEN_ADDRESS '0.0.0.0'
ENV GRPC_SERVER_LISTEN_PORT '9999'
ENV SERVICE_NAME_AS_ROOT 'false'

ENV DATABASE_DIALECT 'mysql'
ENV DATABASE_USERNAME 'root'
ENV DATABASE_PASSWORD 'root'
ENV DATABASE_HOST 'mysql-service'
ENV DATABASE_PORT '3306'
ENV DATABASE_NAME 'be-candle'

ENV REDIS_ENDPOINT 'redis-service:6379'
ENV REDIS_DB '0'
ENV REDIS_POOLSIZE '100'
ENV REDIS_IDLE_TIMEOUT '60000'
ENV REDIS_MAX_IDLE_CONNECTIONS '10'
ENV REDIS_CONNECTION_IDLE_TIMEOUT_MS '60000'
ENV REDIS_PASSWORD '123456'

ENV SERVICE_NAME 'be-candle'
ENV SERVER_ENV 'dev'
ENV PRODUCTION_ENVIRONMENT 'false'
ENV LOG_LEVEL '4'
ENV STACKDRIVER_ENABLED 'false'
ENV GRPC_CONNECT_TIMEOUT_MS '15000'
ENV PROJECT_ID 'paper-trade-chatbot'
ENV FCM_KEY '.'
ENV CDN_URL_PREFIX "https://storage.googleapis.com"
ENV SERVER_SHUTDOWN_GRACE_PERIOD_MS '30000'
ENV GCS_BUCKET_NAME 'paper-trade-chatbot-bucket'

ENV PRODUCT_GRPC_HOST 'be-product-service'
ENV PRODUCT_GRPC_PORT '9999'

RUN apk add --update-cache tzdata
COPY be-candle /be-candle

ENTRYPOINT ["/be-candle"]

