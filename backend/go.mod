module github.com/qoal/file-processor

go 1.24.10

require (
	// ===========================================
	// IMAGE CONVERSION LIBRARIES
	// ===========================================

	// Core image processing and format conversions (JPG, PNG, GIF, BMP, TIFF, WebP)
	github.com/disintegration/imaging v1.6.2

	// Core Web Framework & Dependencies
	github.com/gin-contrib/cors v1.7.2
	github.com/gin-gonic/gin v1.10.1
	github.com/go-redis/redis/v8 v8.11.5
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1

	// PDF creation from text/documents (PDF generation)
	github.com/jung-kurt/gofpdf v1.16.2

	// ===========================================
	// ARCHIVE CONVERSION LIBRARIES
	// ===========================================

	// Multi-format archive support: ZIP, TAR, 7Z, RAR, GZ, BZ2, XZ, LZ4
	github.com/mholt/archiver/v3 v3.5.1

	// ===========================================
	// DOCUMENT CONVERSION LIBRARIES
	// ===========================================

	// Microsoft Office formats (DOCX, XLSX, PPTX, DOC, XLS, PPT)
	github.com/unidoc/unioffice v1.35.0

	// Cryptographic functions (password protection, encryption)
	golang.org/x/crypto v0.25.0

	// Extended image format support (additional codecs)
	golang.org/x/image v0.32.0

	// Database drivers
	gorm.io/driver/postgres v1.5.4
	gorm.io/gorm v1.30.0
)

require (

	// ===========================================
	// ARCHIVE SUPPORT LIBRARIES
	// ===========================================

	// Brotli compression algorithm support
	github.com/andybalholm/brotli v1.0.1 // indirect

	// ===========================================
	// CORE FRAMEWORK DEPENDENCIES
	// ===========================================

	// Web framework components
	github.com/bytedance/sonic v1.11.6 // indirect
	github.com/bytedance/sonic/loader v0.1.1 // indirect

	// Hash function utilities
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/cloudwego/base64x v0.1.4 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect

	// Redis caching support
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect

	// Multi-format compression library
	github.com/dsnet/compress v0.0.2-0.20210315054119-f66993602bf5 // indirect

	// MIME type detection
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect

	// Audio buffer management and PCM data handling
	github.com/go-audio/audio v1.0.0 // indirect

	// RIFF file format support (used by WAV)
	github.com/go-audio/riff v1.0.0 // indirect

	// WAV audio format encoding/decoding
	github.com/go-audio/wav v1.1.0
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.20.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect

	// Snappy compression algorithm
	github.com/golang/snappy v0.0.2 // indirect
	// ===========================================
	// AUDIO CONVERSION LIBRARIES
	// ===========================================

	// MP3 decoding library for audio processing
	github.com/hajimehoshi/go-mp3 v0.3.4

	// Database drivers
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.4.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/json-iterator/go v1.1.12 // indirect

	// Compression libraries
	github.com/klauspost/compress v1.11.4 // indirect
	github.com/klauspost/cpuid/v2 v2.2.7 // indirect
	github.com/klauspost/pgzip v1.2.5 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect

	// RAR archive format decoder
	github.com/nwaples/rardecode v1.1.0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect

	// LZ4 compression algorithm
	github.com/pierrec/lz4/v4 v4.1.2 // indirect

	// ===========================================
	// DOCUMENT SUPPORT LIBRARIES
	// ===========================================

	// Microsoft Office document parsing (OLE format)
	github.com/richardlehane/msoleps v1.0.3 // indirect

	// Testing utilities
	github.com/stretchr/testify v1.11.1 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect

	// XZ compression algorithm
	github.com/ulikunitz/xz v0.5.9 // indirect

	// Additional XZ compression support
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect

	// System libraries
	golang.org/x/arch v0.8.0 // indirect
	golang.org/x/net v0.27.0 // indirect
	golang.org/x/sys v0.22.0 // indirect
	golang.org/x/text v0.30.0 // indirect

	// Protocol buffer support
	google.golang.org/protobuf v1.34.1 // indirect

	// YAML configuration support
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
