package repositories

import (
	"Project/config"
	"Project/models"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestAllSLmethods(t *testing.T) {

	// Set up database connection
	godotenv.Load()
	db, err := config.SetupDB()
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	repo := NewSLRepository(db)
	//viRepo := NewVoucherItemRepository(db)
	//vRepo := NewVoucherRepository(db)

	repo.db.AutoMigrate(&models.SL{})

	var code, title string

	t.Run("SL CRUD", func(t *testing.T) {
		t.Run("Create", func(t *testing.T) {

			t.Run("Should not Add SL with empty code", func(t *testing.T) {
				code, title = "", fmt.Sprintf("%d", time.Now().UnixNano())
				assert.NotNil(t, repo.Create(code, title, false))
			})

			t.Run("Should not Add SL with empty title", func(t *testing.T) {
				code, title = fmt.Sprintf("%d", time.Now().UnixNano()), ""
				assert.NotNil(t, repo.Create(code, title, false))
			})

			t.Run("Should not Add SL with a code that has more that 64 chars", func(t *testing.T) {
				code, title = fmt.Sprintf("%d", time.Now().UnixNano()), "acceptable title"
				for range 65 {
					code += "a"
				}
				assert.NotNil(t, repo.Create(code, title, false))
			})

			t.Run("Should not Add SL with a title that has more that 64 chars", func(t *testing.T) {
				code, title = "acceptable code", fmt.Sprintf("%d", time.Now().UnixNano())
				for range 65 {
					title += "a"
				}
				assert.NotNil(t, repo.Create(code, title, false))
			})

			t.Run("Should be Able to add a SL with code that has 64 persian chars", func(t *testing.T) {
				code, title = stamp("")
				for range 64 - len(code) {
					code += "م"
				}
				assert.Nil(t, repo.Create(code, title, true))

				// Check for correct insertion in database
				SL := &models.SL{}
				repo.db.Model(&models.SL{}).Last(SL)
				assert.Equal(t, SL.Code, code)
				assert.Equal(t, SL.Title, title)
				assert.Equal(t, SL.HasDL, true)
			})

			t.Run("Should be Able to add a SL with title that has 64 persian chars", func(t *testing.T) {
				code, title = stamp("")
				for range 64 - len(title) {
					title += "م"
				}
				assert.Nil(t, repo.Create(code, title, false))

				// Check for correct insertion in database
				SL := &models.SL{}
				repo.db.Model(&models.SL{}).Last(SL)
				assert.Equal(t, SL.Code, code)
				assert.Equal(t, SL.Title, title)
				assert.Equal(t, SL.HasDL, false)

			})

			t.Run("Shouldn't Add a SL with repetative code", func(t *testing.T) {
				code, title = stamp("S001")
				err := repo.Create(code, title, false)
				assert.Nil(t, err)

				err = repo.Create(code, title+"new", false)
				assert.NotNil(t, err)
			})

			t.Run("Shouldn't Add a SL with repetative title", func(t *testing.T) {
				code, title = stamp("S002")
				err := repo.Create(code, title, false)
				assert.Nil(t, err)

				err = repo.Create(code+"new", title, false)
				assert.NotNil(t, err)
			})

		})

		t.Run("Read", func(t *testing.T) {
			SL := &models.SL{}
			if err := repo.db.Model(&models.SL{}).First(SL).Error; err != nil {
				fmt.Println("DB issue, aborting...")
				return
			}

			readSL, err := repo.Read(SL.ID)
			assert.Nil(t, err)
			assert.Equal(t, SL, readSL)
		})

		t.Run("Update", func(t *testing.T) {
			SL := &models.SL{}
			if err := repo.db.Model(&models.SL{}).Last(SL).Error; err != nil {
				fmt.Println("DB issue, aborting...")
				return
			}

			t.Run("Shouldnt Update SL Code to empty string", func(t *testing.T) {
				newSL := &models.SL{
					DL: models.DL{
						Code:  "",
						Title: "Valid" + fmt.Sprintf("%d", time.Now().UnixNano()),
					},
				}

				err := repo.Update(SL.ID, newSL)
				assert.NotNil(t, err)
			})

			t.Run("Shouldnt Update SL Title to empty string", func(t *testing.T) {
				newSL := &models.SL{
					DL: models.DL{
						Code:  "Valid" + fmt.Sprintf("%d", time.Now().UnixNano()),
						Title: "",
					},
				}

				err := repo.Update(SL.ID, newSL)
				assert.NotNil(t, err)
			})

			t.Run("Shouldnt Update SL code to string with more than 64 chars", func(t *testing.T) {
				code, title = stamp("")
				for range 65 - len(code) {
					code += "o"
				}
				newSL := &models.SL{
					DL: models.DL{
						Code:  code,
						Title: title,
					},
				}

				err := repo.Update(SL.ID, newSL)
				assert.NotNil(t, err)
			})

			t.Run("Shouldnt Update SL title to string with more than 64 chars", func(t *testing.T) {
				code, title = stamp("")
				for range 65 - len(title) {
					title += "o"
				}
				newSL := &models.SL{
					DL: models.DL{
						Code:  code,
						Title: title,
					},
				}

				err := repo.Update(SL.ID, newSL)
				assert.NotNil(t, err)
			})

			t.Run("Shouldnt Update SL code to sth repetetive", func(t *testing.T) {
				firstSL := &models.SL{}
				repo.db.Model(&models.SL{}).First(firstSL)

				newSL := &models.SL{
					DL: models.DL{
						Code:  firstSL.Code,
						Title: "sth valid",
					},
				}

				err := repo.Update(SL.ID, newSL)
				assert.NotNil(t, err)
			})

			t.Run("Shouldnt Update SL titel to sth repetetive", func(t *testing.T) {
				firstSL := &models.SL{}
				repo.db.Model(&models.SL{}).First(firstSL)

				newSL := &models.SL{
					DL: models.DL{
						Code:  "sth valid",
						Title: firstSL.Title,
					},
				}

				err := repo.Update(SL.ID, newSL)
				assert.NotNil(t, err)
			})

			t.Run("should update to a valid state", func(t *testing.T) {
				code, title = stamp("S003-Updated")
				newSL := &models.SL{
					DL: models.DL{
						Code:  code,
						Title: title,
					},
				}
				err := repo.Update(SL.ID, newSL)

				assert.Nil(t, err)

				updatedSL := &models.SL{}
				repo.db.Model(&models.SL{}).First(updatedSL, "id = ?", SL.ID)

				assert.Equal(t, code, updatedSL.Code)
				assert.Equal(t, title, updatedSL.Title)

			})

			t.Run("Row Versioning", func(t *testing.T) {
				code, title = stamp("S003-Updated-Twice")
				newSL := &models.SL{
					DL: models.DL{
						Code:  code,
						Title: title,
					},
				}
				err := repo.Update(SL.ID, newSL)
				assert.NotNil(t, err)

				newSL.Version = 1
				err = repo.Update(SL.ID, newSL)
				//fmt.Println(err)
				assert.Nil(t, err)

				code, title = stamp("S003-Updated-Triple")
				newSL = &models.SL{
					DL: models.DL{
						Code:  code,
						Title: title,
					},
				}
				err = repo.Update(SL.ID, newSL)
				assert.NotNil(t, err)

				newSL.Version = 2
				err = repo.Update(SL.ID, newSL)
				assert.Nil(t, err)

				code, title = stamp("S003-Updated-4th")
				newSL = &models.SL{
					DL: models.DL{
						Code:  code,
						Title: title,
					},
				}
				err = repo.Update(SL.ID, newSL)
				assert.NotNil(t, err)

			})

		})

		t.Run("Delete", func(t *testing.T) {

			t.Run("Should Delete the last 4 SL's", func(t *testing.T) {
				sl := &models.SL{}
				if err := repo.db.Model(&models.SL{}).Last(sl).Error; err != nil {
					fmt.Println("DB issue, aborting...", err)
					return
				}
				err := repo.Delete(sl.ID)
				assert.Nil(t, err)
				_, err = repo.Read(sl.ID)
				assert.NotNil(t, err)

				sl = &models.SL{}
				if err := repo.db.Model(&models.SL{}).Last(sl).Error; err != nil {
					fmt.Println("DB issue, aborting...", err)
					return
				}

				err = repo.Delete(sl.ID)
				assert.Nil(t, err)
				_, err = repo.Read(sl.ID)
				assert.NotNil(t, err)

				sl = &models.SL{}
				if err := repo.db.Model(&models.SL{}).Last(sl).Error; err != nil {
					fmt.Println("DB issue, aborting...", err)
					return
				}

				err = repo.Delete(sl.ID)
				assert.Nil(t, err)
				_, err = repo.Read(sl.ID)
				assert.NotNil(t, err)

				sl = &models.SL{}
				if err := repo.db.Model(&models.SL{}).Last(sl).Error; err != nil {
					fmt.Println("DB issue, aborting...", err)
					return
				}

				err = repo.Delete(sl.ID)
				assert.Nil(t, err)
				_, err = repo.Read(sl.ID)
				assert.NotNil(t, err)
			})

		})
	})
}
