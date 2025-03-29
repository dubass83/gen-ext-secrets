## Generate External Secrets from Kubernetes Secrets

This repository provides a utility to generate external secrets from Kubernetes `secrets.data`. It simplifies the process of extracting and transforming Kubernetes secrets into a format compatible with external secret management systems.

### Features

- Parse Kubernetes `secrets.data` files.
- Generate external secret configurations.

### Usage

1. Clone the repository:
    ```bash
    git clone https://github.com/dubass83/gen-ext-secrets.git
    cd gen-ext-secrets
    ```

2. Run the utility:
    ```bash
    cd  /utils/k8s-secrets && chmod +x ./generate-secrets.sh
    ./generate-secrets.sh && rm generate-secrets.sh
    ```

3. Run go program
    Edit constant values in the main.go file with apropriate values.
    ```golang
    const (
	    pathToJsonDir       = "./utils"
	    pathToGenExtSecrets = "./ext-secrets"
	    yamlTemplate        = "./template/ext-secret.yaml.gotmpl"
	    Enviroment          = "devel"       // Set to "prod" for production
	    KubeNamespace       = "default"     // Set the default namespace
	    VaultPath           = "secret/test" // Set the default Vault path
    )
    ```
    Run the program
    ```bash
    go run ./main.go
    ```

4. Output will be saved in the specified format.

### Requirements

- Bash or compatible shell.
- Kubernetes secrets file in YAML format.

### Contributing

Contributions are welcome! Please open an issue or submit a pull request.

### License

This project is licensed under the MIT License. 