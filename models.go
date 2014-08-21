package pressure

import (
	"database/sql"
	"reflect"
	"time"
)

type QuerySet interface {
	Filter() QuerySet
	OrderBy(field string) QuerySet
	Execute() ([]ModelInstance, error)
	Get() (ModelInstance, error)
}

type ModelInstance interface {
	Save() error
	Delete() error
}

type Model interface {
	New() ModelInstance
	Objects() QuerySet
	TableName() string
	Query() ([]ModelInstance, error)
	CreateTable()
	DropTable()
}

type ModelFieldPrimaryKey int
type ModelFieldForeignKey int
type ModelFieldInt int
type ModelFieldString string
type ModelFieldText string
type ModelFieldTime *time.Time
type ModelFieldBlob []byte

type BasicQuerySet struct{}

type ModelField struct {
	dbName string
	dbType reflect.Type
	value  interface{}
}

func (m *ModelField) TypeName() string {
	switch m.value.(type) {
	case ModelFieldPrimaryKey:
		return "int"
	case ModelFieldForeignKey:
		return "int"
	case ModelFieldInt:
		return "int"
	}
	return "undef"
}

type BasicModel struct {
	tableName string
	fields    map[string]*ModelField
	empty     interface{}
}

func (b *BasicModel) New() ModelInstance {
	return &BasicModelInstance{
		baseModel: b,
		isNew:     true,
	}
}
func (b *BasicModel) Objects() QuerySet               { return nil }
func (b *BasicModel) TableName() string               { return b.tableName }
func (b *BasicModel) Query() ([]ModelInstance, error) { return nil, nil }

func (b *BasicModel) DropTable() {}

func (b *BasicModel) CreateTable() {
	var output = "CREATE TABLE " + b.TableName() + " ("
	for n, v := range b.fields {
		output += "\n" + n + " " + v.TypeName()
	}
	// CREATE TABLE table_name
	// (
	// column_name1 data_type(size),
	// )
}

type BasicModelInstance struct {
	baseModel    *BasicModel
	isNew        bool
	filledFields map[string]*ModelField
}

func (b *BasicModelInstance) Delete() error { return nil }
func (b *BasicModelInstance) Save() error   { return nil }

type ModelEngine struct {
	tablePrefix string
	models      map[string]Model
	db          *sql.DB
	*Logger
}

func (s *Server) NewModelEngineFromDatabaseString(drivername string, conn string, tablePrefix string) (*ModelEngine, error) {
	db, err := sql.Open(drivername, conn)
	if err != nil {
		s.LogError("Unable to open database connection.")
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		s.LogError("Cannot connect to database.")
		return nil, err
	}

	a := &ModelEngine{
		tablePrefix: tablePrefix,
		db:          db,
		Logger:      s.Logger,
	}
	return a, nil
}

func (m *ModelEngine) Close() {
	m.db.Close()
}

// Create a Model from a Struct, then Register it
func (m *ModelEngine) RegisterStruct(s interface{}) Model {
	// Get Underlying Types
	t := reflect.ValueOf(&s).Elem()
	typeOfT := t.Type()

	// Create BasicModel from Struct
	b := &BasicModel{
		tableName: typeOfT.Name(),
		empty:     s,
	}

	values := make(map[string]*ModelField, 0)

	// Loop Through Struct Fields
	for i := 0; i < t.NumField(); i++ {
		f := t.Type().Field(i)

		// Create ModelField from this information...
		i := &ModelField{
			dbName: f.Name,
			dbType: f.Type,
		}

		values[f.Name] = i
	}

	b.fields = values
	return b
}

func (m *ModelEngine) QuerySQL(load_object interface{}, command string, args ...interface{}) (int, error) {
	return 0, nil
}

func (m *ModelEngine) ExecuteSQL(command string, args ...interface{}) (int, error) {
	return 0, nil
}

// Add a custom model to the engine
func (m *ModelEngine) RegisterModel(s Model) {
	m.models[s.TableName()] = s
}

func (m *ModelEngine) CreateModels() {}
func (m *ModelEngine) FlushModels()  {}

// // Flush the Database
// if *db_flush {
// 	theServer.DbMap.DropTables()
// }

// // Create the Database
// if *db_create {
// 	// Create Tables
// 	err := theServer.DbMap.CreateTablesIfNotExists()
// 	if err != nil {
// 		fmt.Println("Problem Creating Tables")
// 		fmt.Println(err)
// 		return
// 	}

// 	// The Tracker
// 	fmt.Println("Time to setup the default Trackers")

// 	for {
// 		var t_url, t_address string
// 		fmt.Print("Tracker URL (or 'done' to stop): ")
// 		fmt.Scanln(&t_url)

// 		if t_url == "done" {
// 			break
// 		}

// 		fmt.Print("Tracker Address: ")
// 		fmt.Scanln(&t_address)

// 		t_full := &models.Tracker{URL: t_url, Address: t_address}
// 		theServer.DbMap.Insert(t_full)
// 	}
// 	// Create the User
// 	fmt.Println("Let's create the first user.")

// 	var username, password, first, last string

// 	fmt.Print("Username (no spaces): ")
// 	fmt.Scanln(&username)

// 	fmt.Print("Password: ")
// 	fmt.Scanln(&password)

// 	fmt.Print("First Name: ")
// 	fmt.Scanln(&first)
// 	fmt.Print("Last Name: ")
// 	fmt.Scanln(&last)

// 	newUser := models.CreateUser(username, password, theServer)
// 	newUser.FullName = (first + " " + last)

// 	fmt.Println("New User Address", newUser.Address)
// 	theServer.DbMap.Insert(newUser)
// }
