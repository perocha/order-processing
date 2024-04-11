# order-processing
Microservice to process orders

# Build & Deployment

## Creating local docker image

To create a local docker images for each microservice, you can use the following command:

```bash
docker image build -t order-processing -f .\cmd\order-processing\Dockerfile .
```

# Configuration

## Environment Variables

The configuration is managed using Azure App Configuration. The microservice only requires one environment variable with the connection string to the App Config service:
* App configuration: APPCONFIGURATION_CONNECTION_STRING

In the k8s deployment file, this environment variables needs to be set:

```yaml
        env:
        - name: APPCONFIGURATION_CONNECTION_STRING
          valueFrom:
            secretKeyRef:
              name: appconfiguration
              key: appconfigurationconnectionstring
```

And then, the secret needs to be created in the AKS cluster:

```bash
kubectl create secret generic appconfiguration --from-literal=appconfigurationconnectionstring="<connection string>"
```

# Testing

Send a message to the publisher service, you can use curl:

```bash
curl -X POST -H "Content-Type: application/json" -d "{\"Type\": \"create_order\"}" http://<ipaddress>:80/publish
```

