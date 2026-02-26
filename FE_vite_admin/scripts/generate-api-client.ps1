$ErrorActionPreference = "Stop"

$repoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..\..")).Path
$specPath = Join-Path $repoRoot "docs\FE\openapi.yaml"
$outputPath = Join-Path $repoRoot "FE_vite\src\api\generated"

if (-not (Test-Path $specPath)) {
  throw "OpenAPI spec not found: $specPath"
}

Write-Host "Generating API client from: $specPath"

docker run --rm -v "${repoRoot}:/local" openapitools/openapi-generator-cli generate `
  -i /local/docs/FE/openapi.yaml `
  -g typescript-axios `
  -o /local/FE_vite/src/api/generated `
  --additional-properties=npmName=golf-store-api-client,supportsES6=true,withSeparateModelsAndApi=true,modelPackage=models,apiPackage=services,typescriptThreePlus=true,enumPropertyNaming=original

Write-Host "API client generated at: $outputPath"
