package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	pathToJsonDir       = "../vault-secrets"
	pathToGenExtSecrets = "../ext-secrets"
	yamlTemplate        = "./template/ext-secret.yaml.gotmpl"
	Enviroment          = "devel"       // Set to "prod" for production
	KubeNamespace       = "default"     // Set the default namespace
	VaultPath           = "secret/test" // Set the default Vault path
)

type TemplateData struct {
	Name      string
	Namespace string
	VaultPath string
	JsonData  map[string]any
}

func main() {
	// Set the log level based on the environment
	if Enviroment == "devel" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	// Read the JSON files from the pathToJsonDir directory
	jsonFiles, err := readJsonFiles(pathToJsonDir)
	if err != nil {
		log.Fatal().Err(err).Msg("Error reading JSON files")
	}

	// Generate the ext-secrets.yaml file
	err = generateExtSecretsFile(jsonFiles, pathToGenExtSecrets, yamlTemplate)
	if err != nil {
		log.Fatal().Err(err).Msg("Error generating ext-secrets.yaml file")
	}
}

// readJsonFiles reads JSON files from the specified directory and returns a slice of file names.
// It returns an error if the directory cannot be read or if there are no JSON files.
func readJsonFiles(path string) ([]string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Error().Err(err).Msg("Error reading directory")
		return nil, err
	}

	var jsonFiles []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			log.Info().Msgf("Found JSON file: %s", file.Name())
			// Append the full path of the JSON file to the slice
			jsonFiles = append(jsonFiles, filepath.Join(path, file.Name()))
		}
	}

	if len(jsonFiles) == 0 {
		log.Error().Msg("No JSON files found in directory")
		return nil, fmt.Errorf("no JSON files found in directory: %s", path)
	}

	return jsonFiles, nil
}

// generateExtSecretsFile generates the ext-secrets.yaml file based on the provided JSON files.
// It returns an error if the file cannot be created or written to.
func generateExtSecretsFile(jsonFiles []string, outputPath string, templatePath string) error {
	// Read the template file
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		log.Error().Err(err).Msg("Error reading template file")
		return fmt.Errorf("failed to read template file: %w", err)
	}

	// Create the output directory if it doesn't exist
	err = os.MkdirAll(outputPath, os.ModePerm)
	if err != nil {
		log.Error().Err(err).Msg("Error creating output directory")
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Iterate over the JSON files and generate the ext-secrets.yaml content
	for _, jsonFile := range jsonFiles {
		// Read the JSON file content
		jsonContent, err := os.ReadFile(jsonFile)
		if err != nil {
			log.Error().Err(err).Msgf("Error reading JSON file: %s", jsonFile)
			return fmt.Errorf("failed to read JSON file %s: %w", jsonFile, err)
		}

		// Process the template with the JSON content (this assumes a templating library is used)
		// For simplicity, we'll just append the JSON content to the output file
		outputFile := filepath.Join(outputPath, "ext-secrets.yaml")
		f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Error().Err(err).Msgf("Error opening output file: %s", outputFile)
			return fmt.Errorf("failed to open output file %s: %w", outputFile, err)
		}
		defer f.Close()
		// Use a templating engine to process the template with the JSON content
		tmpl, err := template.New("ext-secrets").Parse(string(templateContent))
		if err != nil {
			log.Error().Err(err).Msg("Error parsing template file")
			return fmt.Errorf("failed to parse template file: %w", err)
		}

		// Populate the template data

		tmplateData := TemplateData{}
		tmplateData.Name = baseName(jsonFile) // Use the base name of the JSON file as the name
		tmplateData.Namespace = KubeNamespace // Set the namespace as needed
		tmplateData.VaultPath = VaultPath

		err = json.Unmarshal(jsonContent, &tmplateData.JsonData)
		if err != nil {
			log.Error().Err(err).Msgf("Error unmarshalling JSON file: %s", jsonFile)
			return fmt.Errorf("failed to unmarshal JSON file %s: %w", jsonFile, err)
		}

		// Write the processed template to the output file
		err = tmpl.Execute(f, tmplateData)
		if err != nil {
			log.Error().Err(err).Msgf("Error executing template for JSON file: %s", jsonFile)
			return fmt.Errorf("failed to execute template for JSON file %s: %w", jsonFile, err)
		}

		log.Info().Msgf("Processed JSON file: %s", jsonFile)
	}

	log.Info().Msg("Successfully generated ext-secrets.yaml file")
	return nil
}

func baseName(path string) string {
	// Get the base name of the file without the extension
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	return strings.TrimSuffix(base, ext)
}
