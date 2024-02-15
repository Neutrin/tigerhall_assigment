package repositiories

import (
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
		Updates(map[string]interface{}{"last_seen": sighting.LastSeenTimeStamp, "latititude": sighting.Lat, "longitude": sighting.Long}).
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
	err = repo.db.Debug().Scopes(Paginate(tigers, &pagParams, repo.db)).Find(&tigers).Error
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
