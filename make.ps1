param(
    [string]$Platform = "all"  # 支持: linux-amd64, linux-arm64, all
)
$originalEnv = @{
    GOOS = $env:GOOS
    GOARCH = $env:GOARCH
}
$projectRoot = $PSScriptRoot
$outputDir = "$projectRoot\bin"
$package = "."  # 修改为你的主包路径

# 目标平台配置
$targets = @()
if ($Platform -eq "all") {
    $targets = @(
        @{ GOOS = "linux"; GOARCH = "amd64"; Suffix = "" },
        @{ GOOS = "linux"; GOARCH = "arm64"; Suffix = "" }
        # 可添加其他平台
    )
} else {
    $parts = $Platform.Split('-')
    $targets += @{ GOOS = $parts[0]; GOARCH = $parts[1]; Suffix = "" }
}

# 创建输出目录
if (-not (Test-Path $outputDir)) { New-Item -ItemType Directory $outputDir | Out-Null }

try {
    foreach ($target in $targets) {
        $env:GOOS = $target.GOOS
        $env:GOARCH = $target.GOARCH
        $outputName = "$(Split-Path $projectRoot -Leaf)_$($env:GOOS)-$($env:GOARCH)$($target.Suffix)"

        Write-Host "Compiling for ${env:GOOS}/${env:GOARCH}..."
        go build -o "$outputDir\$outputName" "$package"

        if ($LASTEXITCODE -ne 0) {
            Write-Error "编译失败! 目标: ${env:GOOS}/${env:GOARCH}"
            exit 1
        }
    }
    Write-Host "✅ 所有目标编译完成! 输出目录: $outputDir"
}
finally {
    # 确保任何情况下都还原环境
    $env:GOOS = $originalEnv.GOOS
    $env:GOARCH = $originalEnv.GOARCH
    $env:CGO_ENABLED = $originalEnv.CGO_ENABLED
    Write-Host "♻️ 环境变量已还原"
}