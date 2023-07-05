FROM go:1.20
LABEL authors="noname"

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY *.go ./
COPY .version/*.go ./.version/
COPY utils/*.go ./utils/
COPY cmd/*.go ./cmd/
COPY scan/*.go ./scan/
COPY rules/*.go ./rules/
COPY utils/*.go ./utils/

# Copy target
COPY target/* ./target/

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /noscan

# Run
CMD ["/noscan"]