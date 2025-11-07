package services

import (
	"fmt"
	"os"
	"strings"

	"github.com/jung-kurt/gofpdf"
	"github.com/qoal/file-processor/models"
	"github.com/unidoc/unioffice/document"
	"github.com/unidoc/unioffice/spreadsheet"
)

func (p *EnhancedDocumentProcessor) convertTextToPDF(input, output string, job *models.ProcessingJob) (string, error) {
	// Read text content
	content, err := os.ReadFile(input)
	if err != nil {
		return "", fmt.Errorf("failed to read text file: %w", err)
	}

	// Create PDF using gofpdf
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set font
	pdf.SetFont("Arial", "", 12)

	// Add text content
	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		pdf.Text(10, 20+float64(i*8), line)
	}

	// Write PDF to output file
	if err := pdf.OutputFileAndClose(output); err != nil {
		return "", fmt.Errorf("failed to write PDF file: %w", err)
	}

	return output, nil
}

func (p *EnhancedDocumentProcessor) convertDocxToText(input, output string, job *models.ProcessingJob) (string, error) {
	// Open DOCX document
	doc, err := document.Open(input)
	if err != nil {
		return "", fmt.Errorf("failed to open DOCX: %w", err)
	}
	defer doc.Close()

	// Extract text content
	var textContent strings.Builder

	// Extract text from paragraphs
	for _, para := range doc.Paragraphs() {
		for _, run := range para.Runs() {
			textContent.WriteString(run.Text())
		}
		textContent.WriteString("\n")
	}

	// Write text to output file
	if err := os.WriteFile(output, []byte(textContent.String()), 0644); err != nil {
		return "", fmt.Errorf("failed to write text file: %w", err)
	}

	return output, nil
}

func (p *EnhancedDocumentProcessor) convertTextToDocx(input, output string, job *models.ProcessingJob) (string, error) {
	// Read text content
	content, err := os.ReadFile(input)
	if err != nil {
		return "", fmt.Errorf("failed to read text file: %w", err)
	}

	// Create new DOCX document
	doc := document.New()
	defer doc.Close()

	// Add text content as paragraphs
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		para := doc.AddParagraph()
		run := para.AddRun()
		run.AddText(line)
	}

	// Save document
	if err := doc.SaveToFile(output); err != nil {
		return "", fmt.Errorf("failed to save DOCX file: %w", err)
	}

	return output, nil
}

func (p *EnhancedDocumentProcessor) convertXlsxToCSV(input, output string, job *models.ProcessingJob) (string, error) {
	// Open XLSX spreadsheet
	ss, err := spreadsheet.Open(input)
	if err != nil {
		return "", fmt.Errorf("failed to open XLSX: %w", err)
	}
	defer ss.Close()

	// Get the first sheet
	sheet := ss.Sheets()[0]

	// Convert to CSV format
	var csvContent strings.Builder

	// Iterate through rows and cells
	for _, row := range sheet.Rows() {
		var cells []string
		for _, cell := range row.Cells() {
			cells = append(cells, cell.GetString())
		}
		csvContent.WriteString(strings.Join(cells, ","))
		csvContent.WriteString("\n")
	}

	// Write CSV to output file
	if err := os.WriteFile(output, []byte(csvContent.String()), 0644); err != nil {
		return "", fmt.Errorf("failed to write CSV file: %w", err)
	}

	return output, nil
}

func (p *EnhancedDocumentProcessor) convertCSVToXlsx(input, output string, job *models.ProcessingJob) (string, error) {
	// Read CSV content
	content, err := os.ReadFile(input)
	if err != nil {
		return "", fmt.Errorf("failed to read CSV file: %w", err)
	}

	// Create new spreadsheet
	ss := spreadsheet.New()
	defer ss.Close()

	// Add a sheet
	sheet := ss.AddSheet()

	// Parse CSV and populate sheet
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		cells := strings.Split(line, ",")
		row := sheet.AddRow()

		for _, cellValue := range cells {
			cell := row.AddCell()
			cell.SetString(cellValue)
		}
	}

	// Save spreadsheet
	if err := ss.SaveToFile(output); err != nil {
		return "", fmt.Errorf("failed to save XLSX file: %w", err)
	}

	return output, nil
}

func (p *EnhancedDocumentProcessor) mergePDFs(inputFiles []string, output string, job *models.ProcessingJob) (string, error) {
	// Create output file
	out, err := os.Create(output)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	// For now, create a simple merged representation
	var mergedContent strings.Builder
	mergedContent.WriteString("Merged PDF document\n")
	mergedContent.WriteString(fmt.Sprintf("Merged from %d files:\n", len(inputFiles)))

	for i, file := range inputFiles {
		mergedContent.WriteString(fmt.Sprintf("File %d: %s\n", i+1, file))
	}

	if err := os.WriteFile(output, []byte(mergedContent.String()), 0644); err != nil {
		return "", fmt.Errorf("failed to write merged PDF: %w", err)
	}

	return output, nil
}
