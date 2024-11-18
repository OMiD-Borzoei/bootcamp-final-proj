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

func stamp(code string) (string, string) {
	code = fmt.Sprintf("%s-%d", code, time.Now().UnixNano())
	return code, code + "-T"
}

func TestAllDLmethods(t *testing.T) {

	// Set up database connection
	godotenv.Load()
	db, err := config.SetupDB()
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	repo := NewDLRepository(db)
	//viRepo := NewVoucherItemRepository(db)
	//vRepo := NewVoucherRepository(db)

	repo.db.AutoMigrate(&models.DL{})

	var code, title string

	t.Run("DL CRUD", func(t *testing.T) {
		t.Run("Create", func(t *testing.T) {

			t.Run("Should not Add DL with empty code", func(t *testing.T) {
				code, title = "", fmt.Sprintf("%d", time.Now().UnixNano())
				assert.NotNil(t, repo.Create(code, title))
			})

			t.Run("Should not Add DL with empty title", func(t *testing.T) {
				code, title = fmt.Sprintf("%d", time.Now().UnixNano()), ""
				assert.NotNil(t, repo.Create(code, title))
			})

			t.Run("Should not Add DL with a code that has more that 64 chars", func(t *testing.T) {
				code, title = fmt.Sprintf("%d", time.Now().UnixNano()), "acceptable title"
				for range 65 {
					code += "a"
				}
				assert.NotNil(t, repo.Create(code, title))
			})

			t.Run("Should not Add DL with a title that has more that 64 chars", func(t *testing.T) {
				code, title = "acceptable code", fmt.Sprintf("%d", time.Now().UnixNano())
				for range 65 {
					title += "a"
				}
				assert.NotNil(t, repo.Create(code, title))
			})

			t.Run("Should be Able to add a DL with code that has 64 persian chars", func(t *testing.T) {
				code, title = stamp("")
				for range 64 - len(code) {
					code += "م"
				}
				assert.Nil(t, repo.Create(code, title))

				// Check for correct insertion in database
				dl := &models.DL{}
				repo.db.Model(&models.DL{}).Last(dl)
				assert.Equal(t, dl.Code, code)
				assert.Equal(t, dl.Title, title)
			})

			t.Run("Should be Able to add a DL with title that has 64 persian chars", func(t *testing.T) {
				code, title = stamp("")
				for range 64 - len(title) {
					title += "م"
				}
				assert.Nil(t, repo.Create(code, title))

				// Check for correct insertion in database
				dl := &models.DL{}
				repo.db.Model(&models.DL{}).Last(dl)
				assert.Equal(t, dl.Code, code)
				assert.Equal(t, dl.Title, title)

			})

			t.Run("Shouldn't Add a DL with repetative code", func(t *testing.T) {
				code, title = stamp("D001")
				err := repo.Create(code, title)
				assert.Nil(t, err)

				err = repo.Create(code, title+"new")
				assert.NotNil(t, err)
			})

			t.Run("Shouldn't Add a DL with repetative title", func(t *testing.T) {
				code, title = stamp("D002")
				err := repo.Create(code, title)
				assert.Nil(t, err)

				err = repo.Create(code+"new", title)
				assert.NotNil(t, err)
			})

		})

		t.Run("Read", func(t *testing.T) {
			dl := &models.DL{}
			if err := repo.db.Model(&models.DL{}).First(dl).Error; err != nil {
				fmt.Println("DB issue, aborting...")
				return
			}

			readdl, err := repo.Read(dl.ID)
			assert.Nil(t, err)
			assert.Equal(t, dl, readdl)
		})

		t.Run("Update", func(t *testing.T) {
			dl := &models.DL{}
			if err := repo.db.Model(&models.DL{}).Last(dl).Error; err != nil {
				fmt.Println("DB issue, aborting...")
				return
			}

			t.Run("Shouldnt Update DL Code to empty string", func(t *testing.T) {
				newdl := &models.DL{
					Code:  "",
					Title: "Valid" + fmt.Sprintf("%d", time.Now().UnixNano()),
				}

				err := repo.Update(dl.ID, newdl)
				assert.NotNil(t, err)
			})

			t.Run("Shouldnt Update DL Title to empty string", func(t *testing.T) {
				newdl := &models.DL{
					Code:  "Valid" + fmt.Sprintf("%d", time.Now().UnixNano()),
					Title: "",
				}

				err := repo.Update(dl.ID, newdl)
				assert.NotNil(t, err)
			})

			t.Run("Shouldnt Update DL code to string with more than 64 chars", func(t *testing.T) {
				code, title = stamp("")
				for range 65 - len(code) {
					code += "o"
				}
				newdl := &models.DL{
					Code:  code,
					Title: title,
				}

				err := repo.Update(dl.ID, newdl)
				assert.NotNil(t, err)
			})

			t.Run("Shouldnt Update DL title to string with more than 64 chars", func(t *testing.T) {
				code, title = stamp("")
				for range 65 - len(title) {
					title += "o"
				}
				newdl := &models.DL{
					Code:  code,
					Title: title,
				}

				err := repo.Update(dl.ID, newdl)
				assert.NotNil(t, err)
			})

			t.Run("Shouldnt Update DL code to sth repetetive", func(t *testing.T) {
				firstdl := &models.DL{}
				repo.db.Model(&models.DL{}).First(firstdl)

				newdl := &models.DL{
					Code:  firstdl.Code,
					Title: "sth valid",
				}

				err := repo.Update(dl.ID, newdl)
				assert.NotNil(t, err)
			})

			t.Run("Shouldnt Update DL titel to sth repetetive", func(t *testing.T) {
				firstdl := &models.DL{}
				repo.db.Model(&models.DL{}).First(firstdl)

				newdl := &models.DL{
					Code:  "sth valid",
					Title: firstdl.Title,
				}

				err := repo.Update(dl.ID, newdl)
				assert.NotNil(t, err)
			})

			t.Run("should update to a valid state", func(t *testing.T) {
				code, title = stamp("D003-Updated")
				newdl := &models.DL{
					Code:  code,
					Title: title,
				}
				err := repo.Update(dl.ID, newdl)
				assert.Nil(t, err)

				updateddl := &models.DL{}
				repo.db.Model(&models.DL{}).First(updateddl, "id = ?", dl.ID)

				assert.Equal(t, code, updateddl.Code)
				assert.Equal(t, title, updateddl.Title)

			})

			t.Run("Row Versioning", func(t *testing.T) {
				code, title = stamp("D003-Updated-Twice")
				newdl := &models.DL{
					Code:  code,
					Title: title,
				}
				err := repo.Update(dl.ID, newdl)
				assert.NotNil(t, err)

				newdl.Version = 1
				err = repo.Update(dl.ID, newdl)
				assert.Nil(t, err)

				code, title = stamp("D003-Updated-Triple")
				newdl = &models.DL{
					Code:  code,
					Title: title,
				}
				err = repo.Update(dl.ID, newdl)
				assert.NotNil(t, err)

				newdl.Version = 2
				err = repo.Update(dl.ID, newdl)
				assert.Nil(t, err)

				code, title = stamp("D003-Updated-4th")
				newdl = &models.DL{
					Code:  code,
					Title: title,
				}
				err = repo.Update(dl.ID, newdl)
				assert.NotNil(t, err)

			})

		})

		t.Run("Delete", func(t *testing.T) {
			dl := &models.DL{}
			if err := repo.db.Model(&models.DL{}).Last(dl).Error; err != nil {
				fmt.Println("DB issue, aborting...")
				return
			}

			t.Run("Shouldn't delete a DL which is referenced in at least 1 voucheritem", func(t *testing.T) {
				//

			})

		})
	})
}
