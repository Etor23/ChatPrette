$projectPath = "c:\Users\brony\Downloads\ChatPrette Back\ChatPrette"
cd $projectPath

# 1. Actualizar internal\dto\user_dto.go
$dtoFile = Get-Content "internal\dto\user_dto.go" -Raw
if ($dtoFile -notmatch "UpdateProfileRequest") {
    $newDto = $dtoFile.TrimEnd() + "`n`n// UpdateProfileRequest para cambiar username o birthdate`ntype UpdateProfileRequest struct {`n`tUsername  string `"json:``"username,omitempty`"`"  // opcional`n`tBirthdate string `"json:``"birthdate,omitempty`"`" // opcional`n`tAvatarURL string `"json:``"avatar_url,omitempty`"`" // opcional`n}"
    Set-Content "internal\dto\user_dto.go" $newDto
    Write-Host "✓ Updated internal\dto\user_dto.go"
}

# 2. Actualizar internal\repos\user_repo.go
$repoFile = Get-Content "internal\repos\user_repo.go" -Raw
if ($repoFile -notmatch "func.*Update.*ctx.*user") {
    $updateMethod = @"

func (r *UserRepo) Update(ctx context.Context, id string, username *string, birthdate *time.Time, avatarURL *string) (*models.User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("ID inválido: %w", err)
	}

	update := bson.M{
		"updated_at": time.Now(),
	}

	// Solo actualizar los campos que se proporcionaron
	if username != nil && *username != "" {
		update["username"] = *username
	}
	if birthdate != nil {
		update["birthdate"] = birthdate
	}
	if avatarURL != nil && *avatarURL != "" {
		update["avatar_url"] = *avatarURL
	}

	opts := options.FindOneAndUpdateOptions{}
	opts.SetReturnDocument(options.After)

	var updatedUser models.User
	err = r.collection.FindOneAndUpdate(ctx, bson.M{"_id": objID}, bson.M{"$set": update}, &opts).Decode(&updatedUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("usuario no encontrado")
		}
		return nil, err
	}

	return &updatedUser, nil
}
"@
    $newRepo = $repoFile.TrimEnd() + $updateMethod
    Set-Content "internal\repos\user_repo.go" $newRepo
    Write-Host "✓ Added Update method to internal\repos\user_repo.go"
}

# 3. Agregar "strings" import a auth_handler.go si no está
$handlerFile = Get-Content "internal\handlers\auth_handler.go" -Raw
if ($handlerFile -notmatch '"strings"') {
    $handlerFile = $handlerFile -replace '(\nimport \([^)]*)"net/http"', "`$1`"net/http`"`n`t`"strings`""
    Set-Content "internal\handlers\auth_handler.go" $handlerFile
    Write-Host "✓ Added strings import to internal\handlers\auth_handler.go"
}

# 4. Agregar UpdateProfile handler al auth_handler.go
if ($handlerFile -notmatch "func.*UpdateProfile") {
    $updateHandler = @"

// PUT /api/auth/me
// Actualiza el perfil del usuario (username, birthdate, avatar_url)
// REQUIERE: Authorization header con Bearer <token>
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	// El middleware ya verificó el token
	userID := auth.GetUserID(c)

	var body dto.UpdateProfileRequest

	// 1. Parsear y validar body
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	// 2. Si cambia username, verificar que no esté en uso (excepto el suyo)
	if body.Username != "" {
		currentUser, err := h.repo.FindById(ctx, userID)
		if err != nil || currentUser == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
			return
		}

		// Si el nuevo username es diferente al actual, verificar disponibilidad
		if body.Username != currentUser.Username {
			taken, err := h.repo.ExistsByUsername(ctx, body.Username)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al verificar username"})
				return
			}
			if taken {
				c.JSON(http.StatusConflict, gin.H{"error": "Username ya en uso"})
				return
			}
		}
	}

	// 3. Parsear birthdate si se proporciona
	var birthPtr *time.Time
	if body.Birthdate != "" {
		// try RFC3339 first, then date-only
		if t, err := time.Parse(time.RFC3339, body.Birthdate); err == nil {
			birthPtr = &t
		} else if t2, err2 := time.Parse("2006-01-02", body.Birthdate); err2 == nil {
			birthPtr = &t2
		}
	}

	// 4. Actualizar en la BD
	username := &body.Username
	if body.Username == "" {
		username = nil
	}

	avatarURL := &body.AvatarURL
	if body.AvatarURL == "" {
		avatarURL = nil
	}

	user, err := h.repo.Update(ctx, userID, username, birthPtr, avatarURL)
	if err != nil {
		if strings.Contains(err.Error(), "no encontrado") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar perfil"})
		}
		return
	}

	// 5. Formatear respuesta
	var birthStr string
	if user.Birthdate != nil {
		birthStr = user.Birthdate.Format(time.RFC3339)
	}
	createdStr := user.CreatedAt.Format(time.RFC3339)

	c.JSON(http.StatusOK, dto.UserResponse{
		ID:        user.ID.Hex(),
		Email:     user.Email,
		Username:  user.Username,
		AvatarURL: user.AvatarURL,
		Birthdate: birthStr,
		CreatedAt: createdStr,
	})
}
"@
    $newHandler = $handlerFile.TrimEnd() + $updateHandler
    Set-Content "internal\handlers\auth_handler.go" $newHandler
    Write-Host "✓ Added UpdateProfile handler to internal\handlers\auth_handler.go"
}

# 5. Actualizar router.go para registrar la ruta
$routerFile = Get-Content "router.go" -Raw
if ($routerFile -notmatch "UpdateProfile") {
    $routerFile = $routerFile -replace '(protectedAuth\.POST\("/refresh", authHandler\.Refresh\))', "`$1`n`t`t`t`tprotectedAuth.PUT(`"/me`", authHandler.UpdateProfile)"
    Set-Content "router.go" $routerFile
    Write-Host "✓ Added UpdateProfile route to router.go"
}

Write-Host "`n✅ Todos los cambios completados"
