package breeds_importer_service

import (
	"context"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const catSheetIndex = 0
const dogSheetIndex = 1

type BreedsImporterService struct {
	client *sheets.Service
}

func NewBreedsImporterService(ctx context.Context, apiKey string) (*BreedsImporterService, error) {
	client, err := sheets.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	return &BreedsImporterService{client: client}, nil
}

func (c *BreedsImporterService) GetSpreadsheet(spreadsheetId string) (*sheets.Spreadsheet, error) {
	resp := c.client.Spreadsheets.Get(spreadsheetId)
	resp.IncludeGridData(true)
	spreadsheet, err := resp.Do()

	if err != nil {
		return nil, err
	}

	return spreadsheet, nil
}

func (c *BreedsImporterService) GetCatNames(spreadsheet *sheets.Spreadsheet) []Row {
	var catRows []Row

	var catsSheet *sheets.Sheet = spreadsheet.Sheets[catSheetIndex]
	for _, row := range catsSheet.Data[0].RowData[1:] {
		if len(row.Values) == 0 {
			continue
		}

		catRows = append(catRows, parseRow(row))
	}

	return catRows
}

func (c *BreedsImporterService) GetDogNames(spreadsheet *sheets.Spreadsheet) []Row {
	var dogRows []Row

	var dogsSheet *sheets.Sheet = spreadsheet.Sheets[dogSheetIndex]
	for _, row := range dogsSheet.Data[0].RowData[1:] {
		if len(row.Values) == 0 {
			continue
		}

		dogRows = append(dogRows, parseRow(row))
	}

	return dogRows
}

type Row struct {
	Breed string
}

func parseRow(row *sheets.RowData) Row {
	return Row{
		Breed: row.Values[1].FormattedValue,
	}
}
