# Chirpstack Simulator

A simulator for Chirpstack, a LoRaWAN Network Server stack.

Adapted from [Chirpstack Simulator](github.com/brocaar/chirpstack-simulator).

## Deployment

### Quick Start

1. **Add your configuration**

   Create a `config.toml` file in the `/manifests` directory.  
   This file should contain your Chirpstack Simulator configuration.  
   For example:

   ```toml
   # /manifests/config.toml
   [chirpstack]
   # Add your ChirpStack configuration here
   ```

2. **Build and Deploy to Kubernetes**

   From the project root, run:

   ```sh
   make deploy
   ```

   This will:
   - Build and push the Docker image (if needed)
   - Apply the Kubernetes manifests using Kustomize
   - Create or update the Kubernetes Secret containing your `config.toml`
   - Redeploy the simulator deployment and trigger a rollout

3. **Verify the Deployment**

   Check that the pods are running:

   ```sh
   kubectl get pods -n chirpstack
   ```

   You can also check logs:

   ```sh
   kubectl logs deployment/chirpstack-simulator-nextevolution -n chirpstack
   ```

---

**Note:**  

- You may need to adjust the image tag or registry in the `Makefile` before deploying.
- The deployment expects the secret `chirpstack-simulator-config` to be mounted, which is automatically handled by Kustomize.
- Make sure your Kubernetes context is set correctly or adjust the `KUBECONFIG` variable in the `Makefile` as needed.
