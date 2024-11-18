package repositories

import (
	"Project/config"
	"Project/models"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func scenario2(db *gorm.DB) (sl1_id, sl2_id, dl_id uint) {
	slRepo := NewSLRepository(db)
	dlRepo := NewDLRepository(db)

	dlRepo.db.AutoMigrate(&models.DL{}, &models.SL{}, &models.Voucher{}, &models.Voucheritem{})

	code, title := stamp("S010")
	sl1_id, _ = slRepo.Create(code, title, true)

	code, title = stamp("S011")
	sl2_id, _ = slRepo.Create(code, title, false)

	code, title = stamp("D010")
	dl_id, _ = dlRepo.Create(code, title)

	return sl1_id, sl2_id, dl_id
}

func revert_scenario2(db *gorm.DB, sl1_id, sl2_id, dl_id uint) {
	slRepo := NewSLRepository(db)
	dlRepo := NewDLRepository(db)

	slRepo.Delete(sl1_id)
	slRepo.Delete(sl2_id)
	dlRepo.Delete(dl_id)
}

func TestAllVmethods(t *testing.T) {
	// Set up database connection
	godotenv.Load()
	db, err := config.SetupDB()
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	repo := NewVoucherRepository(db)

	itemedV := models.Voucher{}
	sl1_id, sl2_id, dl_id := scenario2(db)

	vi1 := models.Voucheritem{
		SLID:   sl1_id,
		DLID:   &dl_id,
		Credit: 200,
		Debit:  0,
	}

	vi2 := models.Voucheritem{
		SLID:   sl2_id,
		DLID:   nil,
		Credit: 0,
		Debit:  200,
	}

	repo.db.AutoMigrate(&models.Voucher{}, &models.DL{}, &models.SL{})

	t.Run("Voucher CRUD", func(t *testing.T) {
		t.Run("Create", func(t *testing.T) {

			t.Run("Should only create a v that has [2, 500] vis", func(t *testing.T) {
				itemedV.Number, _ = stamp("V001")
				_, err := repo.Create(&itemedV)
				assert.Contains(t, err.Error(), "expected between [2, 500]")

				itemedV.Number, _ = stamp("V001")
				itemedV.Items = append(itemedV.Items, vi1)
				_, err = repo.Create(&itemedV)
				assert.Contains(t, err.Error(), "expected between [2, 500]")

				itemedV.Number, _ = stamp("V001")
				itemedV.Items = append(itemedV.Items, vi2)
				_, err = repo.Create(&itemedV)
				assert.Nil(t, err)

				for range 499 {
					itemedV.Items = append(itemedV.Items, vi1)
				}

				itemedV.Number, _ = stamp("V001")
				_, err = repo.Create(&itemedV)
				assert.Contains(t, err.Error(), "expected between [2, 500]")

				itemedV.Items = make([]models.Voucheritem, 0, 2)
				itemedV.Items = append(itemedV.Items, vi1)
				itemedV.Items = append(itemedV.Items, vi2)
			})

			t.Run("Shouldnt create an unbalanced V", func(t *testing.T) {
				itemedV.Number, _ = stamp("V001")
				itemedV.Items[0].Credit += 1

				_, err = repo.Create(&itemedV)
				assert.Contains(t, err.Error(), "unbalance")

				itemedV.Items[0].Credit -= 1
				itemedV.Items[1].Debit += 1

				itemedV.Number, _ = stamp("V001")
				_, err = repo.Create(&itemedV)
				assert.Contains(t, err.Error(), "unbalance")

				itemedV.Items[1].Debit -= 1
			})

			t.Run("Shouldnt create a V that holds a vi with both debit credit > 0", func(t *testing.T) {
				itemedV.Number, _ = stamp("V001")
				itemedV.Items[0].Credit += 1
				itemedV.Items[0].Debit += 1

				_, err = repo.Create(&itemedV)
				assert.Equal(t, err.Error(), "one of the credit or debit must be 0")

				itemedV.Items[0].Credit -= 1
				itemedV.Items[0].Debit -= 1
			})

			t.Run("Shouldnt create a V with empty number", func(t *testing.T) {
				itemedV.Number, _ = stamp("V001")
				itemedV.Number = ""
				_, err := repo.Create(&itemedV)
				assert.Contains(t, err.Error(), "cannot be nil")
			})

			t.Run("Shouldnt create a V with a number that has more than 64 chars", func(t *testing.T) {
				itemedV.Number, _ = stamp("V001")
				itemedV.Number = ""
				for range 65 - len(itemedV.Number) {
					itemedV.Number += "a"
				}
				_, err := repo.Create(&itemedV)
				assert.Equal(t, err.Error(), "Number Must not have more than 64 characters")
			})

			t.Run("Shouldnt create a V with a repetetive number", func(t *testing.T) {
				itemedV.Number, _ = stamp("V002")

				_, err := repo.Create(&itemedV)
				assert.Nil(t, err)

				_, err = repo.Create(&itemedV)
				assert.Contains(t, err.Error(), "duplicate key")
			})

			t.Run("Should create a V with a number that has 64 Persian chars", func(t *testing.T) {
				itemedV.Number, _ = stamp("V003")
				for range 64 - len(itemedV.Number) {
					itemedV.Number += "ک"
				}
				_, err := repo.Create(&itemedV)
				assert.Nil(t, err)
			})

			t.Run("shouldnt create a v when its vi doesnt provide a valid sl", func(t *testing.T) {
				itemedV.Number, _ = stamp("V004")
				itemedV.Items[1].SLID = 0

				_, err := repo.Create(&itemedV)
				assert.Equal(t, err.Error(), "SL with ID 0 does not exist")

				itemedV.Items[1].SLID = vi2.SLID
			})

			t.Run("shouldnt create a v when its vi has a sl.hasdl=True but no valid dl is provided", func(t *testing.T) {
				itemedV.Number, _ = stamp("V004")
				itemedV.Items[0].DLID = nil

				_, err := repo.Create(&itemedV)
				assert.Contains(t, err.Error(), "sl must have a dl")

				var x uint = 0
				itemedV.Items[0].DLID = &x

				_, err = repo.Create(&itemedV)
				assert.Equal(t, err.Error(), "DL with ID 0 does not exist")

				itemedV.Items[0].DLID = vi1.DLID
			})

			t.Run("shouldnt create a v when its vi has a sl.hasdl=false but a dl is provided", func(t *testing.T) {
				itemedV.Number, _ = stamp("V004")
				var x uint = 0
				itemedV.Items[1].DLID = &x

				_, err := repo.Create(&itemedV)
				assert.Contains(t, err.Error(), "dl_id is provided")

				itemedV.Items[1].DLID = vi2.DLID
			})

		})

		t.Run("Read", func(t *testing.T) {
			v := &models.Voucher{}
			if err := repo.db.Model(&models.Voucher{}).Preload("Items").First(v).Error; err != nil {
				fmt.Println("DB issue, aborting...")
				return
			}

			readv, err := repo.Read(v.ID)
			assert.Nil(t, err)
			assert.Equal(t, v, readv)
		})

		t.Run("Update", func(t *testing.T) {
			firstV := models.Voucher{}
			repo.db.Model(&models.Voucher{}).Last(&firstV)
			itemedV.ID = uint(int(firstV.ID))

			t.Run("Shouldnt Update to an invalid state", func(t *testing.T) {
				t.Run("vis < 2 || vis > 500", func(t *testing.T) {
					last1 := &models.Voucheritem{}
					repo.db.Model(&models.Voucheritem{}).Last(&last1)
					itemedV.Items = []models.Voucheritem{
						{
							ID: uint(last1.ID),
						},

						{
							ID: uint(last1.ID - 1),
						},
					}
					itemedV.Number, _ = stamp("V001-Updated")

					err := repo.Update(itemedV.ID, &itemedV)
					if !assert.Contains(t, err.Error(), "expected between [2, 500]") {
						os.Exit(0)
					}

				})

				t.Run("unbalanced voucher", func(t *testing.T) {
					itemedV.Items = []models.Voucheritem{
						{
							ID:     uint(9),
							SLID:   vi1.SLID,
							DLID:   vi1.DLID,
							Debit:  100,
							Credit: 0,
						},

						{
							ID:     uint(10),
							SLID:   vi2.SLID,
							DLID:   vi2.DLID,
							Debit:  0,
							Credit: 100,
						},
					}
					itemedV.Items[1].Credit += 1
					itemedV.Number, _ = stamp("V001-Updated")
					err := repo.Update(itemedV.ID, &itemedV)
					assert.Contains(t, err.Error(), "unbalance")
				})

				t.Run("debit > 0 && credit > 0 ", func(t *testing.T) {
					itemedV.Items[1].Debit += 1
					itemedV.Number, _ = stamp("V001-Updated")
					err := repo.Update(itemedV.ID, &itemedV)
					assert.Equal(t, err.Error(), "one of the credit or debit must be 0")
					itemedV.Items[1].Debit -= 1
					itemedV.Items[1].Credit -= 1
				})

				t.Run("empty number", func(t *testing.T) {
					itemedV.Number = ""
					err := repo.Update(itemedV.ID, &itemedV)
					assert.Contains(t, err.Error(), "cannot be nil")
				})

				t.Run("number has more than 64 chars", func(t *testing.T) {
					itemedV.Number = ""
					for range 65 {
						itemedV.Number += "a"
					}
					err := repo.Update(itemedV.ID, &itemedV)
					assert.Equal(t, err.Error(), "Number Must not have more than 64 characters")
				})

				t.Run("number has more than 64 chars", func(t *testing.T) {
					itemedV.Number = ""
					for range 65 {
						itemedV.Number += "a"
					}
					err := repo.Update(itemedV.ID, &itemedV)
					assert.Equal(t, err.Error(), "Number Must not have more than 64 characters")
				})

				t.Run("repetetive number", func(t *testing.T) {
					fetchedV := &models.Voucher{}
					repo.db.Model(&models.Voucher{}).First(fetchedV)

					itemedV.Number = fetchedV.Number
					err := repo.Update(itemedV.ID, &itemedV)

					assert.Contains(t, err.Error(), "duplicate key")
				})

				t.Run("no valid sl provided", func(t *testing.T) {
					itemedV.Number, _ = stamp("V001-Updated")

					itemedV.Items[0].SLID = 0

					err := repo.Update(itemedV.ID, &itemedV)
					assert.Contains(t, err.Error(), "SL with ID 0 does not exist")
				})

				t.Run("sl.haddl = True but no dl provided", func(t *testing.T) {
					itemedV.Number, _ = stamp("V001-Updated")

					itemedV.Items[0].SLID = vi1.SLID
					itemedV.Items[0].DLID = nil

					err := repo.Update(itemedV.ID, &itemedV)
					assert.Contains(t, err.Error(), "sl must have a dl")

					var x uint = 0
					itemedV.Items[0].DLID = &x

					_, err = repo.Create(&itemedV)
					assert.Equal(t, err.Error(), "DL with ID 0 does not exist")

					itemedV.Items[0].DLID = vi1.DLID
				})

				t.Run("sl.haddl = false but a dl provided", func(t *testing.T) {
					var x uint = 0
					itemedV.Items[1].DLID = &x

					_, err := repo.Create(&itemedV)
					assert.Contains(t, err.Error(), "dl_id is provided")

					itemedV.Items[1].DLID = vi2.DLID

				})

			})

			t.Run("Should Update to a valid state", func(t *testing.T) {
				err := repo.Update(itemedV.ID, &itemedV)
				assert.Nil(t, err)
			})

			t.Run("Row Versioning", func(t *testing.T) {
				itemedV.Items[0].Debit = 50
				itemedV.Items[1].Credit = 50

				itemedV.Number, _ = stamp("V001-Updated-Twice")
				itemedV.Version = 1
				err := repo.Update(itemedV.ID, &itemedV)
				assert.Nil(t, err)

				itemedV.Number, _ = stamp("V001-Updated-Triple")
				itemedV.Version = 2
				err = repo.Update(itemedV.ID, &itemedV)
				assert.Nil(t, err)

				itemedV.Number, _ = stamp("V001-Updated-Triple")
				err = repo.Update(itemedV.ID, &itemedV)
				assert.NotNil(t, err)

			})
		})

		t.Run("Delete", func(t *testing.T) {
			t.Run("Should Delete the last 3 Vs", func(t *testing.T) {
				v := &models.Voucher{}
				repo.db.Model(&models.Voucher{}).Last(v)
				assert.Nil(t, repo.Delete(v.ID))
				_, err := repo.Read(v.ID)
				assert.Equal(t, err, gorm.ErrRecordNotFound)

				v = &models.Voucher{}
				repo.db.Model(&models.Voucher{}).Last(v)
				assert.Nil(t, repo.Delete(v.ID))
				_, err = repo.Read(v.ID)
				assert.Equal(t, err, gorm.ErrRecordNotFound)

				v = &models.Voucher{}
				repo.db.Model(&models.Voucher{}).Last(v)
				assert.Nil(t, repo.Delete(v.ID))
				_, err = repo.Read(v.ID)
				assert.Equal(t, err, gorm.ErrRecordNotFound)
			})
		})
	})

	revert_scenario2(db, sl1_id, sl2_id, dl_id)
}
