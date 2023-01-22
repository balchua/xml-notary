/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"strconv"

	"github.com/balchua/xml-notary/pkg/certmgr"
	"github.com/balchua/xml-notary/pkg/notary"
	"github.com/beevik/etree"
	"github.com/gofiber/fiber/v2"
	dsig "github.com/russellhaering/goxmldsig"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// serveCmd represents the serve command
var (
	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Starts the server",
		Long:  `Starts the web server to allow sending xml payloads to sign or verify.`,
		Run:   start,
	}

	port     int
	certFile string
	keyFile  string
	no       *notary.Notary
	ks       dsig.X509KeyStore
)

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().IntVarP(&port, "port", "p", 5000, "configures the port used by the application.")
	serveCmd.Flags().StringVarP(&certFile, "cert", "c", "", "The location of the certificate to use.")
	serveCmd.Flags().StringVarP(&keyFile, "key", "k", "", "The location of the private key to use.")
}

func start(cmd *cobra.Command, args []string) {

	zap.L().Info("list all arguments",
		// Structured context as strongly typed Field values.
		zap.Int("port", port),
		zap.String("certFile", certFile),
		zap.String("keyFile", keyFile),
	)

	var err error

	ks, err = certmgr.New(certFile, keyFile)
	if err != nil {
		zap.L().Fatal("unable to initialize the application ", zap.Error(err))
	}
	no, err = notary.New(ks)

	app := fiber.New(fiber.Config{
		Immutable:     true,
		AppName:       "Xml Notary",
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
	})

	app.Post("/api/sign", func(c *fiber.Ctx) error {
		doc := etree.NewDocument()
		doc.ReadFromBytes(c.Body())
		signedElement, err := no.SignEnvelope(doc.Root())
		if err != nil {
			zap.L().Error("unable to sign document", zap.Error(err))
			return fiber.NewError(fiber.StatusBadRequest, "unable to sign document")
		}

		signedDocument := etree.NewDocument()
		signedDocument.SetRoot(signedElement)
		b, err := signedDocument.WriteToBytes()
		if err != nil {
			zap.L().Error("unable to serialize to bytes", zap.Error(err))
			return fiber.NewError(fiber.StatusBadRequest, "unable to serialize to bytes")
		}
		return c.SendString(string(b))
	})

	app.Listen(":" + strconv.Itoa(port))
}
