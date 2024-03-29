package users

import (
	"cookvs/model"
	"database/sql"
	"fmt"
	"math/rand"
)

type SqlRepo struct {
	DB *sql.DB
}

func RandomImagePath() string {
	min := 0
	max := 5
	img := rand.Intn(max-min) + min
	return fmt.Sprintf("http://localhost:3000/assets/images/%v.png", img)
}

func (r *SqlRepo) Insert(user model.User) error {
	statement, err := r.DB.Prepare(`INSERT INTO ` + "`users`" + ` (
		image,
		email,
		password,
		nickname,
		borndate ) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("failed to prepare insert user: %w", err)
	}
	_, err = statement.Exec(user.Image, user.Email, user.Password, user.NickName, user.BornDate)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)

	}

	return nil
}

func (r *SqlRepo) CheckEmail(user model.User) (*model.User, error) {
	model_user := model.User{}
	rows, err := r.DB.Query(`SELECT * FROM`+"`users`"+`
	WHERE email = ?`, user.Email)
	if err != nil {
		return &model_user, fmt.Errorf("failed to prepare find user: %w", err)
	}

	if rows.Next() {
		rows.Scan(&model_user.UserID, &model_user.Image, &model_user.Email, &model_user.Password, &model_user.NickName, &model_user.BornDate)
		return &model_user, nil
	}

	return nil, nil

}

func (r *SqlRepo) Login(user model.User) (*model.User, error) {
	model_user := model.User{}
	rows, err := r.DB.Query(`SELECT * FROM`+"`users`"+`
	WHERE email = ? AND password =?`, user.Email, user.Password)
	if err != nil {
		return &model_user, fmt.Errorf("failed to prepare find user: %w", err)
	}

	if rows.Next() {
		rows.Scan(&model_user.UserID, &model_user.Image, &model_user.Email, &model_user.Password, &model_user.NickName, &model_user.BornDate)
		return &model_user, nil
	}

	return nil, nil

}

func (r *SqlRepo) FindById(user_id uint) (model.User, error) {
	model_user := model.User{}
	rows, err := r.DB.Query(`SELECT * FROM `+"`users`"+` WHERE user_id = ?`, user_id)
	if err != nil {
		return model_user, fmt.Errorf("failed to prepare find user: %w", err)
	}

	for rows.Next() {
		rows.Scan(&model_user.UserID, &model_user.Image, &model_user.Email, &model_user.Password, &model_user.NickName, &model_user.BornDate)
	}

	return model_user, nil
}

func (r *SqlRepo) DeleteById(user_id uint) error {
	statement, err := r.DB.Prepare(`DELETE FROM ` + "`users`" + ` WHERE user_id = ?`)
	if err != nil {
		return fmt.Errorf("failed to prepare delete user: %w", err)
	}

	_, err = statement.Exec(user_id)
	if err != nil {
		return fmt.Errorf("failed to prepare delete EXEC user: %w", err)
	}
	return nil

}

func (r *SqlRepo) Update(user_id uint, model_user model.User) error {

	statement, err := r.DB.Prepare(`UPDATE ` + "`users`" + ` SET
		image = ?,
        nickname = ?,
		WHERE user_id = ?`)
	if err != nil {
		return fmt.Errorf("failed to prepare update user: %w", err)
	}
	_, err = statement.Exec(model_user.Image, model_user.NickName, model_user.BornDate, user_id)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *SqlRepo) FindAll() ([]model.User, error) {

	array := []model.User{}

	rows, err := r.DB.Query(`SELECT * FROM ` + "`users`" + ``)
	if err != nil {
		return []model.User{}, fmt.Errorf("failed to prepare FindAll user: %w", err)
	}

	for rows.Next() {
		model_user := model.User{}
		rows.Scan(&model_user.UserID, &model_user.Image, &model_user.Email, &model_user.Password, &model_user.NickName, &model_user.BornDate)
		array = append(array, model_user)
	}
	return array, nil
}

func (r *SqlRepo) InsertRecipe(recipe model.Recipe) error {
	statement, err := r.DB.Prepare(`INSERT INTO ` + "`recipes`" + ` (
		user_id,
		name,
		image,
		description,
		products,
		category,
		tag) VALUES (?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("failed to prepare insert recipe: %w", err)
	}
	_, err = statement.Exec(recipe.UserID, recipe.Name, recipe.Image, recipe.MarshalDescription(), recipe.MarshalProducts(), recipe.Category, recipe.Tag)
	if err != nil {
		return fmt.Errorf("failed to insert recipe: %w", err)

	}

	return nil
}

func (r *SqlRepo) FindAllRecipe() ([]model.Recipe, error) {

	array := []model.Recipe{}

	rows, err := r.DB.Query(`SELECT * FROM ` + "`recipes`" + `ORDER BY name`)
	if err != nil {
		return []model.Recipe{}, fmt.Errorf("failed to prepare FindAll recipe: %w", err)
	}

	for rows.Next() {
		model_recipe := model.Recipe{}
		products := ""
		description := ""
		rows.Scan(&model_recipe.RecipeID, &model_recipe.UserID, &model_recipe.Name, &model_recipe.Image, &description, &products, &model_recipe.Category, &model_recipe.Tag)
		model_recipe.UnmarshalSteps(description)
		model_recipe.UnmarshalProducts(products)
		array = append(array, model_recipe)
	}
	return array, nil
}

func (r *SqlRepo) FindByName(recipe model.Recipe) ([]model.Recipe, error) {
	array := []model.Recipe{}
	var stroka string
	if recipe.Name != "" {
		stroka = `SELECT * FROM recipes WHERE name LIKE '` + recipe.Name + `%';`
	} else {
		stroka = `SELECT * FROM recipes;`
	}
	fmt.Println(stroka)
	rows, err := r.DB.Query(stroka)
	if err != nil {
		return []model.Recipe{}, fmt.Errorf("failed to prepare find recipe: %w", err)
	}

	for rows.Next() {
		model_recipe := model.Recipe{}
		description := ""
		products := ""
		rows.Scan(&model_recipe.RecipeID, &model_recipe.UserID, &model_recipe.Name, &model_recipe.Image, &description, &products, &model_recipe.Category, &model_recipe.Tag)
		model_recipe.UnmarshalSteps(description)
		model_recipe.UnmarshalProducts(products)
		array = append(array, model_recipe)
	}

	return array, nil

}
func (r *SqlRepo) RecipeByCategory(recipe model.Recipe) ([]model.Recipe, error) {
	array := []model.Recipe{}
	var stroka = `SELECT * FROM recipes WHERE category = '` + recipe.Category + `';`
	fmt.Println(stroka)
	rows, err := r.DB.Query(stroka)
	if err != nil {
		return []model.Recipe{}, fmt.Errorf("failed to prepare find recipe: %w", err)
	}

	for rows.Next() {
		model_recipe := model.Recipe{}
		description := ""
		products := ""
		rows.Scan(&model_recipe.RecipeID, &model_recipe.UserID, &model_recipe.Name, &model_recipe.Image, &description, &products, &model_recipe.Category, &model_recipe.Tag)
		model_recipe.UnmarshalSteps(description)
		model_recipe.UnmarshalProducts(products)
		array = append(array, model_recipe)
	}

	return array, nil
}

func (r *SqlRepo) RecipeByTag(recipe model.Recipe) ([]model.Recipe, error) {
	array := []model.Recipe{}
	var stroka = `SELECT * FROM recipes WHERE tag = '` + recipe.Tag + `';`
	fmt.Println(stroka)
	rows, err := r.DB.Query(stroka)
	if err != nil {
		return []model.Recipe{}, fmt.Errorf("failed to prepare find recipe: %w", err)
	}

	for rows.Next() {
		model_recipe := model.Recipe{}
		description := ""
		products := ""
		rows.Scan(&model_recipe.RecipeID, &model_recipe.UserID, &model_recipe.Name, &model_recipe.Image, &description, &products, &model_recipe.Category, &model_recipe.Tag)
		model_recipe.UnmarshalSteps(description)
		model_recipe.UnmarshalProducts(products)
		array = append(array, model_recipe)
	}

	return array, nil
}

func (r *SqlRepo) RecipeById(recipe_id uint) (model.Recipe, error) {
	model_recipe := model.Recipe{}
	rows, err := r.DB.Query(`SELECT * FROM `+"`recipes`"+` WHERE recipe_id = ?`, recipe_id)
	if err != nil {
		return model_recipe, fmt.Errorf("failed to prepare find user: %w", err)
	}

	for rows.Next() {
		var description string
		var products string
		err := rows.Scan(
			&model_recipe.RecipeID,
			&model_recipe.UserID,
			&model_recipe.Name,
			&model_recipe.Image,
			&description,
			&products,
			&model_recipe.Category,
			&model_recipe.Tag,
		)
		if err != nil {
			return model_recipe, fmt.Errorf("failed to scan row: %w", err)
		}

		model_recipe.UnmarshalSteps(description)
		model_recipe.UnmarshalProducts(products)
	}
	return model_recipe, nil
}

func (r *SqlRepo) RecipeByUserId(user_id uint) ([]model.Recipe, error) {
	array := []model.Recipe{}

	rows, err := r.DB.Query(`SELECT * FROM `+"`recipes`"+` WHERE user_id = ? ORDER BY name`, user_id)
	if err != nil {
		return []model.Recipe{}, fmt.Errorf("failed to prepare FindAll recipe: %w", err)
	}

	for rows.Next() {
		model_recipe := model.Recipe{}
		products := ""
		description := ""
		rows.Scan(&model_recipe.RecipeID, &model_recipe.UserID, &model_recipe.Name, &model_recipe.Image, &description, &products, &model_recipe.Category, &model_recipe.Tag)
		model_recipe.UnmarshalSteps(description)
		model_recipe.UnmarshalProducts(products)
		array = append(array, model_recipe)
	}
	return array, nil
}
