package controllers

import (
	"fiber-api/config"
	"fiber-api/models"
	"fiber-api/utils"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

func CreateTask(c fiber.Ctx) error {
	var task models.Task
	if err := c.Bind().Body(&task); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	userID := c.Locals("user_id")

	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	task.CreatedByID = c.Locals("user_id").(uint)
	task.UpdatedByID = c.Locals("user_id").(uint)

	if err := utils.Validate.Struct(task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid task data. " + err.Error(),
		})
	}

	if err := config.DB.Create(&task).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to save task",
		})
	}

	// Preload relasi CreatedBy (opsional)
	config.DB.Preload("CreatedBy").First(&task, task.ID)
	config.DB.Preload("UpdatedBy").First(&task, task.ID)

	// Balikin response ke client
	return c.Status(fiber.StatusCreated).JSON(task)

}
func GetTasks(c fiber.Ctx) error {
	var tasks []models.Task

	showDeleted := c.Query("show_deleted")

	db := config.DB
	if showDeleted == "true" {
		db = db.Unscoped()
	}

	result := db.Preload("CreatedBy").
		Preload("UpdatedBy").
		Preload("DeletedBy").
		Find(&tasks)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not fetch tasks",
		})
	}
	return c.JSON(tasks)
}
func GetTask(c fiber.Ctx) error {
	taskID := c.Params("id")
	var task models.Task

	showDeleted := c.Query("show_deleted")

	db := config.DB
	if showDeleted == "true" {
		db = db.Unscoped()
	}

	result := db.
		Preload("CreatedBy").
		Preload("UpdatedBy").
		Preload("DeletedBy").
		First(&task, taskID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Task not found",
		})
	}
	return c.JSON(task)
}
func UpdateTask(c fiber.Ctx) error {
	taskID := c.Params("id")
	var task models.Task
	if err := c.Bind().Body(&task); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}
	task.UpdatedByID = c.Locals("user_id").(uint)
	if err := utils.Validate.Struct(task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid task data. " + err.Error(),
		})
	}
	if err := config.DB.Model(&task).Where("id = ?", taskID).Updates(task).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update task",
		})
	}
	// Preload relasi CreatedBy (opsional)
	config.DB.Preload("CreatedBy").First(&task, taskID)
	config.DB.Preload("UpdatedBy").First(&task, taskID)
	// Balikin response ke client
	return c.Status(fiber.StatusOK).JSON(task)
}
func DeleteTask(c fiber.Ctx) error {
	taskID := c.Params("id")
	var task models.Task
	if err := config.DB.First(&task, taskID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Task not found",
		})
	}
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}
	userIDValue := c.Locals("user_id").(uint)
	task.DeletedByID = &userIDValue

	if err := config.DB.Save(&task).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update deleted_by_id",
		})
	}

	if err := config.DB.Delete(&task).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to delete task",
		})
	}
	// Preload relasi CreatedBy (opsional)
	config.DB.Unscoped().Preload("CreatedBy").First(&task, taskID)
	config.DB.Unscoped().Preload("UpdatedBy").First(&task, taskID)
	config.DB.Unscoped().Preload("DeletedBy").First(&task, taskID)
	// Balikin response ke client
	return c.Status(fiber.StatusOK).JSON(task)
}

func GetUserTasks(c fiber.Ctx) error {
	paramID := c.Params("id")
	var userID uint

	if paramID == "" {
		// Kalau gak ada :id → pake user dari JWT
		uid := c.Locals("user_id")
		if uid == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}
		userID = uid.(uint)
	} else {
		// Kalau ada :id → pake itu
		parsedID, err := strconv.Atoi(paramID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
		}
		userID = uint(parsedID)
	}

	showDeleted := c.Query("show_deleted")

	db := config.DB
	if showDeleted == "true" {
		db = db.Unscoped()
	}

	var tasks []models.Task
	result := db.
		Preload("CreatedBy").
		Preload("UpdatedBy").
		Preload("DeletedBy").
		Where("created_by_id = ?", userID).Find(&tasks)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not fetch tasks",
		})
	}
	return c.JSON(tasks)
}
func GetTaskByStatus(c fiber.Ctx) error {
	status := c.Params("status")
	userID := c.Locals("user_id").(uint)

	validStatuses := map[string]bool{
		"pending": true,
		"done":    true,
	}

	if !validStatuses[status] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid status",
		})
	}

	var tasks []models.Task

	if err := config.DB.
		Preload("CreatedBy").
		Preload("UpdatedBy").
		Where("created_by_id = ? AND status = ?", userID, status).
		Find(&tasks).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch tasks by status",
		})
	}

	return c.JSON(tasks)
}
