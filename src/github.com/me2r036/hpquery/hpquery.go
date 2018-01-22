package main

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/net/proxy"
	"log"
	"net"
	"os"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	sshHost string = "hk.php9.cc"
	sshPort string = "22"
	sshUser string = "shinetechchina"
	sshPass string = "password"

	proxyHost string = "192.168.1.157"
	proxyPort string = "7070"

	inDBHost     string = "hpolsprod.clmlyplvpqds.ap-south-1.rds.amazonaws.com"
	inDBPort     string = "3306"
	inDBUsername string = "renjinfeng"
	inDBPassword string = "Shinetech@2017@hp"
	inDBName     string = "hpolsproduction"

	goSite string = "golang.org:80"
)

var dialer string = "tcp"
var myDialer MyDialer

func chooseDialer() {
	_, err := net.DialTimeout("tcp", goSite, time.Second)
	if err != nil {
		dialer = "mydialer"
		myDialer = &ProxyDialer{}
		_, err := myDialer.Dial(goSite)
		if err != nil {
			myDialer = &SSHDialer{}
			fmt.Println("SSH client running...")
		} else {
			fmt.Println("Proxy client running...")
		}
	} else {
		fmt.Println("Direct connecting...")
	}
}

type MyDialer interface {
	Dial(addr string) (net.Conn, error)
}

type ProxyDialer struct{}

func (proxyDiler *ProxyDialer) Dial(addr string) (net.Conn, error) {
	dialer, err := proxy.SOCKS5("tcp",
		proxyHost+":"+proxyPort, nil, proxy.Direct)
	conn, err := dialer.Dial("tcp", addr)

	return conn, err
}

type SSHDialer struct{}

func (sshDialer *SSHDialer) Dial(addr string) (net.Conn, error) {
	conn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		log.Fatal(err)
	}
	agentClient := agent.NewClient(conn)

	sshConfig := &ssh.ClientConfig{
		User: sshUser,
		Auth: []ssh.AuthMethod{},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	if agentClient != nil {
		sshConfig.Auth = append(sshConfig.Auth,
			ssh.PublicKeysCallback(agentClient.Signers))
	}

	if sshPass != "" {
		sshConfig.Auth = append(sshConfig.Auth,
			ssh.PasswordCallback(func() (string, error) {
				return sshPass, nil
			}))
	}

	sshClient, err := ssh.Dial("tcp",
		fmt.Sprintf("%s:%s", sshHost, sshPort),
		sshConfig)
	if err != nil {
		log.Fatal(err)
	}

	return sshClient.Dial("tcp", addr)
}

func getDBString(c string) string {
	if c == "in" {
		// return "root:root@tcp(127.0.0.1:3306)/demo_magento2_20171222"
		return inDBUsername + ":" + inDBPassword +
			"@" + dialer +
			"(" + inDBHost + ":" + inDBPort + ")/" + inDBName
	}

	return ""
}

func getDB(c string) *sql.DB {
	if myDialer != nil {
		mysql.RegisterDial(dialer, myDialer.Dial)
	}
	db, err := sql.Open("mysql", getDBString(c))

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("DB connection error!")
	}
	return db
}

func showBundleSkus(country string, skus []string) {
	db := getDB(country)

	var sku string
	s := strings.Join(skus, ",")
	q := `SELECT DISTINCT sku FROM
		(SELECT value, row_id FROM catalog_product_entity_varchar
		   WHERE attribute_id = 170 AND FIND_IN_SET(value, ?)) s
		LEFT JOIN catalog_product_entity c
		ON s.row_id = c.row_id
	      ORDER BY sku`
	stmt, _ := db.Prepare(q)
	rows, err := stmt.Query(s)
	defer stmt.Close()

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println()
	m := "Related bundle products:"
	fmt.Println(m)
	fmt.Println(strings.Repeat("-", utf8.RuneCountInString(m)))
	fmt.Println()

	c := 0
	for rows.Next() {
		err = rows.Scan(&sku)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(sku)
		c++
	}
	fmt.Println()
	fmt.Println("Total records:", c)
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please specify sku list, separated by comma.\n")
		fmt.Println("usage: hpquery sku1[,sku2,...]")
		os.Exit(0)
	}

	chooseDialer()
	showBundleSkus("in", os.Args[1:])
}
