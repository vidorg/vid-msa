package controller

import (
	"context"
	"fmt"
	"strconv"
	"vid-msa/rpc/proto"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
)

// grpc add  "/add/:a/:b" return a + b
func Add(c *fiber.Ctx) error {
	conn, err := grpc.Dial("localhost:4040", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := proto.NewAddServiceClient(conn)

	a, err := strconv.ParseUint(c.Params("a"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid argument A",
		})
	}
	b, err := strconv.ParseUint(c.Params("b"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid argument B",
		})
	}
	req := &proto.Request{A: int64(a), B: int64(b)}
	if res, err := client.Add(context.Background(), req); err == nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"result": fmt.Sprint(res.Result),
		})
	}
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": err.Error(),
	})
}

// grpc Multi  "/add/:a/:b" return a * b
func Multi(c *fiber.Ctx) error {
	conn, err := grpc.Dial("localhost:4040", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := proto.NewAddServiceClient(conn)

	a, err := strconv.ParseUint(c.Params("a"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid argument A",
		})
	}
	b, err := strconv.ParseUint(c.Params("b"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid argument B",
		})
	}
	req := &proto.Request{A: int64(a), B: int64(b)}
	if res, err := client.Multiply(context.Background(), req); err == nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"result": fmt.Sprint(res.Result),
		})
	}
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": err.Error(),
	})
}
