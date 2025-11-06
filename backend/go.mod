module github.com/qoal/file-processor

go 1.24.10

require (
	// Core Web Framework & Dependencies
	github.com/gin-contrib/cors v1.7.2
	github.com/gin-gonic/gin v1.10.1
	github.com/go-redis/redis/v8 v8.11.5
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	golang.org/x/crypto v0.23.0
	gorm.io/driver/postgres v1.5.4
	gorm.io/gorm v1.30.0
	
	// ===========================================
	// DOCUMENT CONVERSION LIBRARIES
	// ===========================================
	
	// Microsoft Office formats (DOCX, XLSX, PPTX)
	github.com/unidoc/unioffice v1.35.0
	
	// PDF processing and text extraction
	github.com/gen2brain/go-fitz v1.24.13
	
	// PDF creation from text/documents
	github.com/jung-kurt/gofpdf v1.16.2
	
	// ===========================================
	// IMAGE CONVERSION LIBRARIES
	// ===========================================
	
	// Core image processing and format conversions
	github.com/disintegration/imaging v1.6.2
	golang.org/x/image v0.32.0
	
	// ===========================================
	// VIDEO CONVERSION LIBRARIES
	// ===========================================
	
	// FFmpeg Go bindings for video format conversions
	// Supports: MP4↔AVI, MOV↔MP4, MP4↔MKV, etc.
	github.com/u2takey/ffmpeg-go v0.5.0
	
	// ===========================================
	// AUDIO CONVERSION LIBRARIES
	// ===========================================
	
	// MP3 decoding and encoding
	github.com/hajimehoshi/go-mp3 v0.3.4
	
	// WAV format support
	github.com/go-audio/wav v1.1.0
	
	// FLAC format support
	github.com/go-audio/flac v1.0.0
	
	// AAC format support
	github.com/go-audio/aac v1.0.0
	
	// Audio processing primitives (used by all audio formats)
	github.com/go-audio/audio v1.0.0
	
	// Bit-level I/O operations for audio data
	github.com/icza/bitio v1.1.0
	
	// ===========================================
	// ARCHIVE CONVERSION LIBRARIES
	// ===========================================
	
	// Multi-format archive support: ZIP, TAR, 7Z, RAR, GZ
	github.com/mholt/archiver/v3 v3.5.1
)

require (
	github.com/bytedance/sonic v1.11.6 // indirect
	github.com/bytedance/sonic/loader v0.1.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/cloudwego/base64x v0.1.4 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.20.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.4.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.7 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	golang.org/x/arch v0.8.0 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	google.golang.org/protobuf v1.34.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)