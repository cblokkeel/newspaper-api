# Use the official Milvus image as a base
FROM milvusdb/milvus:v2.4.0

# Create necessary directories
RUN mkdir -p /var/lib/milvus /milvus/configs

# Copy the configuration file into the container
COPY config/embed.yaml /milvus/configs/embedEtcd.yaml

# Expose necessary ports
EXPOSE 19530 9091 2379

# Health check
HEALTHCHECK --interval=30s --timeout=20s --start-period=90s --retries=3 CMD curl -f http://localhost:9091/healthz || exit 1

# Command to run Milvus standalone
CMD ["milvus", "run", "standalone"]
