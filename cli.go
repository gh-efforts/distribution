package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/liushuochen/gotable"
	"github.com/urfave/cli/v2"
)

var userManager = &cli.Command{
	Name:  "user",
	Usage: "user manager",
	Subcommands: []*cli.Command{
		userView,
		userUpdate,
		userDelete,
	},
}

var dataSetManager = &cli.Command{
	Name:  "dataset",
	Usage: "dataset manager",
	Subcommands: []*cli.Command{
		datasetView,
		datasetUpdate,
		datasetDelete,
		datasetGet,
	},
}
var pieceManager = &cli.Command{
	Name:  "piece",
	Usage: "piece manager",
	Subcommands: []*cli.Command{
		pieceUpdate,
		pieceDelete,
		pieceView,
	},
}

var userView = &cli.Command{
	Name:  "view",
	Usage: "view all users",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "json",
			Usage: "use json output",
		},
	},
	Action: func(ctx *cli.Context) error {
		users := NewUsers()
		err := users.ReadUsersFromFile()
		if err != nil {
			return err
		}
		if ctx.Bool("json") {
			data, err := users.View()
			if err != nil {
				return err
			}
			fmt.Println(data)
			return nil
		}

		table, err := gotable.Create("org", "sps")
		if err != nil {
			return err
		}
		for _, user := range users.List {
			data, _ := json.Marshal(user.Sps)
			table.AddRow([]string{user.Org, string(data)})
		}
		fmt.Println(table)

		return nil
	},
}

var userUpdate = &cli.Command{
	Name:  "add",
	Usage: "add user",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "sp",
			Value:    "",
			Usage:    "specify sp list",
			Required: false,
			Aliases:  []string{"s"},
		},
		&cli.StringFlag{
			Name:     "org",
			Value:    "",
			Usage:    "specify org name",
			Required: true,
			Aliases:  []string{"u"},
		},
		&cli.BoolFlag{
			Name:  "force",
			Value: false,
			Usage: "force update user,cover",
		},
	},
	Action: func(ctx *cli.Context) error {
		users := NewUsers()
		err := users.ReadUsersFromFile()
		if err != nil {
			return err
		}

		user := new(User)
		user.Org = ctx.String("org")
		user.Sps = strings.Split(strings.TrimSpace(ctx.String("sp")), ",")

		if ok := users.Get(user.Org); ok != nil {
			if !ctx.Bool("force") {
				return fmt.Errorf("already exist org %s, if want to update, please add --force\n", user.Org)
			} else {
				users.Update(user)
			}
		} else {
			users.Add(user)
		}

		err = users.WriteUsersToFile()
		if err != nil {
			return err
		}
		fmt.Printf("add user %s success!\n", user.Org)
		return nil

	},
}

var userDelete = &cli.Command{
	Name:  "delete",
	Usage: "delete a org",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "org",
			Value:    "",
			Usage:    "specify org",
			Required: true,
			Aliases:  []string{"u"},
		},
		&cli.BoolFlag{
			Name:  "really-do-it",
			Usage: "must be specified for the action to take effect",
		},
	},
	Action: func(ctx *cli.Context) error {
		org := ctx.String("org")
		if !ctx.Bool("really-do-it") {
			return fmt.Errorf("--really-do-it must be specified for this action to have an effect; you have been warned")
		}

		users := NewUsers()
		err := users.ReadUsersFromFile()
		if err != nil {
			return err
		}

		if ok := users.Delete(org); ok {
			err = users.WriteUsersToFile()
			if err != nil {
				return err
			}
			fmt.Printf("delete %s success!\n", org)
		} else {
			fmt.Printf("delete %s failed!!\n", org)
		}

		return nil

	},
}

var datasetView = &cli.Command{
	Name:  "view",
	Usage: "view all datasets",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "json",
			Usage: "use json output",
		},
	},
	Action: func(ctx *cli.Context) error {
		datasets := NewDataSets()

		err := datasets.ReadDataSetsFromFile()
		if err != nil {
			return err
		}
		if ctx.Bool("json") {
			data, err := datasets.View()
			if err != nil {
				return err
			}
			fmt.Println(data)
			return nil
		}
		table, err := gotable.Create("dataSetName", "duplicate", "spSum", "pieceSum", "pieceSize(TiB)", "carSize(TiB)")
		if err != nil {
			return err
		}
		for _, dataSet := range datasets.List {
			var pieceSize, carSize int64
			var spSum, pieceSum int
			for _, piece := range dataSet.Pieces {
				pieceSize += piece.PieceSize
				carSize += piece.CarSize
				spSum = len(piece.SpInfos)
			}
			pieceSum = len(dataSet.Pieces)

			table.AddRow([]string{dataSet.DataSetName, strconv.Itoa(dataSet.Duplicate), strconv.Itoa(spSum), strconv.Itoa(pieceSum), strconv.FormatFloat(float64(pieceSize)/(1<<40), 'f', -1, 64), strconv.FormatFloat(float64(carSize)/(1<<40), 'f', -1, 64)})

		}

		fmt.Println(table)
		return nil
	},
}
var datasetUpdate = &cli.Command{
	Name:  "add",
	Usage: "add a dataset",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Usage:    "specify dataSet name",
			Required: true,
			Aliases:  []string{"n"},
		},
		&cli.IntFlag{
			Name:     "duplicate",
			Usage:    "specify dataSet duplicate",
			Required: true,
			Aliases:  []string{"d"},
		},
		&cli.StringFlag{
			Name:     "filepath",
			Usage:    "specify dataSet filepath. must include pieceCid,pieceSize,carSize",
			Required: true,
			Aliases:  []string{"f"},
		},
		&cli.BoolFlag{
			Name:  "force",
			Value: false,
			Usage: "force update dataset,cover",
		},
	},
	Action: func(ctx *cli.Context) error {
		dataSetName := ctx.String("name")
		filePath := ctx.String("filepath")
		duplicate := ctx.Int("duplicate")

		datasets := NewDataSets()
		dataSet := NewDataSet()

		f, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for {
			if scanner.Scan() {
				piece := new(Piece)
				var err = json.Unmarshal(scanner.Bytes(), piece)
				if err != nil {
					return err
				}
				dataSet.Pieces = append(dataSet.Pieces, piece)
			} else {
				break
			}
		}

		dataSet.DataSetName = dataSetName
		dataSet.Duplicate = duplicate

		err = datasets.ReadDataSetsFromFile()
		if err != nil {
			return err
		}
		if ok := datasets.GetDataset(dataSetName); ok != nil {
			if !ctx.Bool("force") {
				return fmt.Errorf("already exist dateset %s, if want to update, please add --force", dataSet.DataSetName)
			} else {
				datasets.UpdateDataSet(dataSet)
			}

		} else {
			datasets.AddDataSet(dataSet)
		}

		err = datasets.WriteDataSetsToFile()
		if err != nil {
			return err
		}
		fmt.Printf("add dataset %s success!\n", dataSet.DataSetName)
		return nil
	},
}

var datasetDelete = &cli.Command{
	Name:  "delete",
	Usage: "delete a dataset",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Usage:    "specify dataSet name",
			Required: true,
			Aliases:  []string{"n"},
		},
		&cli.BoolFlag{
			Name:  "really-do-it",
			Usage: "must be specified for the action to take effect",
		},
	},
	Action: func(ctx *cli.Context) error {
		dataSetName := ctx.String("name")
		if !ctx.Bool("really-do-it") {
			return fmt.Errorf("--really-do-it must be specified for this action to have an effect; you have been warned")
		}

		datasets := NewDataSets()
		err := datasets.ReadDataSetsFromFile()
		if err != nil {
			return err
		}
		if ok := datasets.DeleteDataSet(dataSetName); ok {
			err = datasets.WriteDataSetsToFile()
			if err != nil {
				return err
			}
			fmt.Printf("delete dataset %s success!\n", dataSetName)
		} else {
			fmt.Printf("delete dataset %s failed!!!\n", dataSetName)
		}

		return nil
	},
}

var datasetGet = &cli.Command{
	Name:  "get",
	Usage: "get the download link for the dataset",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Usage:    "specify dataSet name",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "sp",
			Usage:    "specify a sp",
			Required: true,
		},
		&cli.Float64Flag{
			Name:     "size",
			Usage:    "specify total pieceSize(TiB)",
			Required: true,
		},
		&cli.IntFlag{
			Name:  "duplicate",
			Usage: "specify dataset duplicate, default use dataSet.Duplicate",
		},
		&cli.Int64Flag{
			Name:  "repeat",
			Usage: "specify dataset repeat",
			Value: 0,
		},
		&cli.StringFlag{
			Name:    "prefix",
			Usage:   "specify url prefix",
			EnvVars: []string{"DIST_PREFIX"},
		},
		&cli.StringFlag{
			Name:    "suffix",
			Usage:   "specify url suffix",
			EnvVars: []string{"DIST_SUFFIX"},
			Value:   ".car",
		},
		&cli.BoolFlag{
			Name:  "really-do-it",
			Usage: "must be specified for the action to take effect",
		},
	},
	Action: func(ctx *cli.Context) error {
		dataSetName := ctx.String("name")
		sp := ctx.String("sp")
		size := int64(ctx.Float64("size") * (1 << 40))
		if ctx.IsSet("duplicate") {
			duplicate = ctx.Int("duplicate")
		}
		repeat = ctx.Int("repeat")
		prefix := ctx.String("prefix")
		suffix := ctx.String("suffix")

		datasets := NewDataSets()
		err := datasets.ReadDataSetsFromFile()
		if err != nil {
			return err
		}

		users := NewUsers()
		err = users.ReadUsersFromFile()
		if err != nil {
			return err
		}
		sps := users.GetSps(sp)

		if len(sps) == 0 {
			return fmt.Errorf("%s does not belong to any organization, please add user sp first", sp)
		}

		dataSet, pieceSize, carSize := datasets.GetDataset(dataSetName).GetSize(sp, size, sps)

		for _, piece := range dataSet.Pieces {
			fmt.Printf("%s%s%s\n", prefix, piece.PieceCid, suffix)
		}
		fmt.Printf("total pieceSize:%v, total carSize: %v, missing pieceSize:%v\n", float64(pieceSize)/(1<<40), float64(carSize)/(1<<40), float64(size-pieceSize)/(1<<40))

		if !ctx.Bool("really-do-it") {
			return fmt.Errorf("--really-do-it must be specified for this action to have an effect; you have been warned")
		} else {
			err = datasets.WriteDataSetsToFile()
			if err != nil {
				return err
			}
		}
		return nil
	},
}

var pieceUpdate = &cli.Command{
	Name:  "add",
	Usage: "add piece",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Usage:    "specify dataSet name",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "pieceCid",
			Usage:    "specify pieceCid",
			Required: true,
		},
		&cli.Int64Flag{
			Name:     "pieceSize",
			Usage:    "specify pieceSize",
			Required: true,
		},
		&cli.Int64Flag{
			Name:     "carSize",
			Usage:    "specify carSize",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "sps",
			Usage:    "specify sps. f01001,f01002",
			Required: true,
		},
		&cli.BoolFlag{
			Name:  "force",
			Value: false,
			Usage: "force update piece,cover",
		},
	},
	Action: func(ctx *cli.Context) error {
		dataSetName := ctx.String("name")
		pieceCid := ctx.String("pieceCid")
		pieceSize := ctx.Int64("pieceSize")
		carSize := ctx.Int64("carSize")
		sps := strings.Split(strings.TrimSpace(ctx.String("sps")), ",")

		dataSets := NewDataSets()
		err := dataSets.ReadDataSetsFromFile()
		if err != nil {
			return err
		}

		piece := new(Piece)
		piece.PieceCid = pieceCid
		piece.PieceSize = pieceSize
		piece.CarSize = carSize
		for _, sp := range sps {
			spInfo := new(SpInfo)
			spInfo.Sp = sp
			spInfo.Num = 1
			piece.SpInfos = append(piece.SpInfos, spInfo)
		}

		dataSet := dataSets.GetDataset(dataSetName)

		if ok := dataSet.Get(piece.PieceCid); ok != nil {
			if !ctx.Bool("force") {
				return fmt.Errorf("already exist piece %s, if want to update, please add --force\n", piece.PieceCid)
			} else {
				dataSet.Update(piece)
			}
		} else {
			dataSet.Add(piece)
		}

		err = dataSets.WriteDataSetsToFile()
		if err != nil {
			return err
		}
		fmt.Printf("add piece %s success!\n", piece.PieceCid)
		return nil

	},
}

var pieceDelete = &cli.Command{
	Name:  "delete",
	Usage: "delete piece",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Usage:    "specify dataSet name",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "pieceCid",
			Usage:    "specify pieceCid",
			Required: true,
		},
		&cli.BoolFlag{
			Name:  "really-do-it",
			Usage: "must be specified for the action to take effect",
		},
	},
	Action: func(ctx *cli.Context) error {
		dataSetName := ctx.String("name")
		pieceCid := ctx.String("pieceCid")
		if !ctx.Bool("really-do-it") {
			return fmt.Errorf("--really-do-it must be specified for this action to have an effect; you have been warned")
		}

		dataSets := NewDataSets()
		err := dataSets.ReadDataSetsFromFile()
		if err != nil {
			return err
		}

		dataSet := dataSets.GetDataset(dataSetName)

		if ok := dataSet.Delete(pieceCid); ok {
			err = dataSets.WriteDataSetsToFile()
			if err != nil {
				return err
			}
			fmt.Printf("delete piece %s success!\n", dataSetName)
		} else {
			fmt.Printf("delete piece %s failed!!!\n", dataSetName)
		}

		return nil
	},
}
var pieceView = &cli.Command{
	Name:  "view",
	Usage: "view  pieces",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Usage:    "specify dataSet name",
			Required: true,
		},
		&cli.BoolFlag{
			Name:  "json",
			Usage: "use json output",
		},
	},
	Action: func(ctx *cli.Context) error {
		dataSetName := ctx.String("name")
		dataSets := NewDataSets()
		err := dataSets.ReadDataSetsFromFile()
		if err != nil {
			return err
		}

		dataSet := dataSets.GetDataset(dataSetName)

		if ctx.Bool("json") {
			data, err := json.MarshalIndent(dataSet, "", "    ")
			if err != nil {
				return err
			}
			fmt.Println(string(data))
			return nil
		}

		table, err := gotable.Create("pieceCid", "pieceSize(GiB)", "carSize(GiB)", "sps")
		if err != nil {
			return err
		}
		for _, piece := range dataSet.Pieces {
			var sps []string
			for _, sp := range piece.SpInfos {
				sps = append(sps, sp.Sp)
			}
			sp, err := json.Marshal(sps)
			if err != nil {
				return err
			}
			table.AddRow([]string{piece.PieceCid, strconv.FormatFloat(float64(piece.PieceSize)/(1<<30), 'f', -1, 64), strconv.FormatFloat(float64(piece.CarSize)/(1<<30), 'f', -1, 64), string(sp)})

		}

		fmt.Println(table)
		return nil
	},
}
