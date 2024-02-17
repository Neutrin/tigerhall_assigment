package repositiories

import (
	"errors"
	"fmt"
	"log"

	"github.com/nitin/tigerhall/core/internal/config"
	"github.com/nitin/tigerhall/core/internal/model"
	repositiories "github.com/nitin/tigerhall/core/internal/repositiories"
	"github.com/nitin/tigerhall/core/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresqlRepo struct {
	db *gorm.DB
}

func NewPostgresqlUserRepo() (repositiories.UserRepo, error) {
	//lazy intialisation
	db, err := intiDB()
	if err != nil {
		return nil, err
	}
	return &PostgresqlRepo{db: db}, nil
}

func NewPostgresqlTigerRepo() (repositiories.TigerRepo, error) {
	//lazy intialisation
	db, err := intiDB()
	if err != nil {
		return nil, err
	}
	return &PostgresqlRepo{db: db}, nil
}

/*
@Validates : user exists or not
@Returns : user email or error
@Does : generates password hash for user
*/
func (repo *PostgresqlRepo) Create(user model.User) (string, error) {
	hash, err := utils.GenerateHashPassword(user.Password)
	if err != nil {
		return "", err
	}
	user.Password = hash
	repo.db.Create(&user)
	return user.Email, err
}

func (repo *PostgresqlRepo) UserExists(email string) bool {
	var user model.User
	repo.db.Where("email = ?", email).First(&user)
	return (user.ID != 0)
}

func (repo *PostgresqlRepo) User(email string) model.User {
	var user model.User
	repo.db.Where("email = ?", email).First(&user)
	return user

}
func (repo *PostgresqlRepo) CreateTiger(tiger model.Tiger, params ...interface{}) (int, error) {
	var (
		err       error
		imagePath string
	)
	tx := repo.db.Begin()
	err = tx.Create(&tiger).Error
	if err != nil {
		tx.Rollback()
		return -1, err
	}
	if len(params) == 1 {
		imagePath, _ = params[0].(string)
	}
	sightings := model.TigerSightings{
		TigerId:           int(tiger.ID),
		LastSeenTimeStamp: tiger.LastSeenTimeStamp,
		Lat:               tiger.Lat,
		Long:              tiger.Long,
		ImagePath:         imagePath,
		CreatedBy:         tiger.CreatedBy,
	}

	err = tx.Create(&sightings).Error
	if err != nil {
		tx.Rollback()
		return -1, err
	}
	tx.Commit()
	return int(tiger.ID), err

}

func (repo *PostgresqlRepo) CreateTigerSighting(sighting model.TigerSightings) (int, error) {
	var err error
	tx := repo.db.Begin()
	err = tx.Debug().Create(&sighting).Error
	if err != nil {
		tx.Rollback()
		return -1, err
	}
	result := tx.Debug().Model(&model.Tiger{}).Where(" id = ?", sighting.TigerId).
		Updates(map[string]interface{}{"last_seen": sighting.LastSeenTimeStamp,
			"latititude": sighting.Lat, "longitude": sighting.Long, "created_by": sighting.CreatedBy}).
		RowsAffected
	if result == 0 {
		tx.Rollback()
		return -1, fmt.Errorf(" failed to update ")
	}
	tx.Commit()
	return int(sighting.ID), err

}

func (repo *PostgresqlRepo) ListAllTigers(pagParams repositiories.Pagination) (*repositiories.Pagination, error) {
	var (
		tigers []*model.Tiger
		err    error
	)
	err = repo.db.Debug().Scopes(Paginate(tigers, &pagParams, repo.db, map[string]interface{}{})).Find(&tigers).Error
	pagParams.Rows = tigers
	return &pagParams, err
}

func (repo *PostgresqlRepo) TigerById(tigerId int) (model.Tiger, error) {
	var (
		tiger model.Tiger
		err   error
	)
	foundRows := repo.db.Where("id = ? ", tigerId).First(&tiger).RowsAffected
	//TODO return relevant error types
	if foundRows == 0 {
		err = fmt.Errorf(" tiger = %d not found", tigerId)
	}
	return tiger, err
}

/*
ListSightings : Fetches all sightings of tiger from database
@Error ignored record not found
*/
func (repo *PostgresqlRepo) ListSightings(tigerId int, pagParams repositiories.Pagination) (
	*repositiories.Pagination, error) {
	var (
		sightings = make([]*model.TigerSightings, 0)
		err       error
	)
	err = repo.db.Debug().Scopes(Paginate(sightings, &pagParams, repo.db, map[string]interface{}{"tiger_id": tigerId})).
		Model(&model.TigerSightings{}).Preload("Tiger").Where("tiger_id = ?", tigerId).Find(&sightings).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return &pagParams, err
	}
	pagParams.Rows = sightings
	log.Println(" ******** came over here in list sighting ******")
	return &pagParams, err
}

func userSightings(tigerId int, db *gorm.DB) ([]uint, error) {
	var (
		userId    = make([]uint, 0)
		err       error
		sightings = make([]*model.TigerSightings, 0)
	)
	err = db.Debug().Model(&model.TigerSightings{}).Where("tiger_id = ?", tigerId).Find(&sightings).Error
	for _, curSighting := range sightings {
		userId = append(userId, curSighting.CreatedBy)
	}
	return userId, err
}

func tigersAll(db *gorm.DB) ([]int, error) {
	var (
		tigers   []*model.Tiger
		tigerIds []int
		err      error
	)
	err = db.Debug().Model(&model.Tiger{}).Select("id").Find(&tigers).Error
	for _, curTiger := range tigers {
		tigerIds = append(tigerIds, int(curTiger.ID))
	}
	return tigerIds, err
}

// UNEXPORTED METHODS
func intiDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.DBHost, config.DBUser, config.DBPassword, config.DBName, config.DBPort, "disable")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	//for creation of model table
	if err := db.AutoMigrate(&model.User{}, &model.Tiger{}, &model.TigerSightings{}); err != nil {
		return nil, err
	}

	log.Println("************Migrated database****************")
	return db, nil
}

func (repo *PostgresqlRepo) InitTigerCreateMap() map[int]map[int]struct{} {
	var (
		tigerCreatedByMap = make(map[int]map[int]struct{})
		err               error
	)
	if tigers, err := tigersAll(repo.db); err == nil {
		for _, curTiger := range tigers {
			if _, exists := tigerCreatedByMap[curTiger]; !exists {
				tigerCreatedByMap[curTiger] = make(map[int]struct{})
			}
			if users, err := userSightings(curTiger, repo.db); err == nil {
				for _, curUser := range users {
					tigerCreatedByMap[curTiger][int(curUser)] = struct{}{}
				}

			}
		}
	}
	if err != nil {
		log.Printf(" failed = %+v\n", err)
	}
	return tigerCreatedByMap
}
