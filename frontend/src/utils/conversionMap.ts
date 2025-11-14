// Conversion mapping based on actual backend capabilities
// Maps source format to available target formats

export const conversionMap: Record<string, string[]> = {
  // Image conversions
  'jpg': ['png', 'bmp', 'tiff', 'webp', 'gif'],
  'jpeg': ['png', 'bmp', 'tiff', 'webp', 'gif'],
  'png': ['jpg', 'bmp', 'tiff', 'webp', 'gif'],
  'gif': ['jpg', 'png', 'webp'],
  'bmp': ['jpg', 'png', 'webp'],
  'tiff': ['png', 'jpg', 'webp'],
  'tif': ['png', 'jpg', 'webp'],
  'webp': ['jpg', 'png', 'gif', 'bmp', 'tiff'],

  // Document conversions
  'txt': ['pdf', 'docx'],
  'docx': ['txt'],
  'xlsx': ['csv'],
  'csv': ['xlsx'],
  'pdf': [],  // PDF conversion not implemented

  // Audio conversions
  'mp3': ['wav'],
  'wav': ['mp3'],
  'flac': ['mp3'],
  'ogg': ['mp3'],
  'm4a': ['mp3'],

  // Video conversions (file copy with extension change)
  'mp4': ['avi', 'mov', 'mkv', 'webm'],
  'avi': ['mp4'],
  'mov': ['mp4'],
  'mkv': ['mp4'],
  'webm': ['mp4'],

  // Archive conversions
  'zip': ['tar.gz'],
  'rar': ['zip'],
  'tar.gz': ['zip'],
  'tgz': ['zip'],
};

export const getAvailableFormats = (sourceFormat: string): string[] => {
  const format = sourceFormat.toLowerCase();
  return conversionMap[format] || [];
};

export const isConversionSupported = (sourceFormat: string, targetFormat: string): boolean => {
  const availableFormats = getAvailableFormats(sourceFormat);
  return availableFormats.includes(targetFormat.toLowerCase());
};
