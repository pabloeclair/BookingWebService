package admin

import (
	"context"
	"cu_coworking_book/go/internal/db"
	"cu_coworking_book/go/internal/pb"
	"cu_coworking_book/go/internal/utils"
	"errors"
	"log"
	"strings"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AdminService struct {
	pb.UnimplementedAdminServiceServer
	mu sync.RWMutex
}

// Записывает логи применения всех хэндлеров. Если какой-то из хэндлеров запустился неудачно, то сообщает о типе и сообщении ошибки.
func LogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("Admin server: %s", info.FullMethod)

	resp, err := handler(ctx, req)

	if err != nil {
		log.Printf("Error: %v", err)
	}

	return resp, err
}

// Создает нового пользователя с указанными администратором полями.
func (s *AdminService) CreateUser(ctx context.Context, req *pb.CreateRequestAdmin) (*pb.Empty, error) {
	// - Если пользователь с указанной почтой уже существует, то вернется ошибка AlreadyExists.

	// - Если почта администратора не найдена — NotFound.
	// - Если пароль администратора не совпал — Unauthenticated.

	// При любых других ошибках – Internal.
	// Показатель успеха — отсутствие ошибки.

	s.mu.RLock()
	_, err := utils.ComparePassword(ctx, req.AdminEmail, req.AdminPassword, true)
	s.mu.RUnlock()

	if err != nil {
		return nil, err
	}

	user := db.User{
		Email:      req.Email,
		FirstName:  req.FirstName,
		SecondName: req.SecondName,
		Patronymic: *req.Patronymic,
		Password:   req.Password,
		Role:       pb.Role_USER.String(),
	}

	s.mu.Lock()
	_, err = db.CreateUser(ctx, user)
	s.mu.Unlock()

	if err != nil {
		if errors.Is(err, db.ErrBadRequest) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return nil, nil
}

// Возвращает список всех пользователей, удовлетворяющих заданному ключу.
func (s *AdminService) SortUser(ctx context.Context, req *pb.GetRequestAdmin) (*pb.GetResponseArray, error) {
	// - Если ни один пользователь с указанной почтой не найден, вернется ошибка NotFound.

	// - Если почта администратора не найдена — NotFound.
	// - Если пароль администратора не совпал — Unauthenticated.

	// При любых других ошибках – Internal.
	// Показатель успеха — информация о всех найденных пользователях.

	s.mu.RLock()
	_, err := utils.ComparePassword(ctx, req.AdminEmail, req.AdminPassword, true)
	s.mu.RUnlock()
	if err != nil {
		return nil, err
	}

	s.mu.RLock()
	res, err := db.GetUserByKey(ctx, &req.SortBy, req.SortKey)
	s.mu.RUnlock()
	if err != nil {
		return nil, utils.CompareErrAndErrNotFound(err)
	}

	var users []*pb.GetResponse
	for r := range res {
		users = append(users, utils.ParseToResult(res[r]))
	}

	return &pb.GetResponseArray{Users: users}, nil
}

// Возвращает информацию о пользователе по id.
func (s *AdminService) GetUserById(ctx context.Context, req *pb.Email) (*pb.GetResponse, error) {
	// - Если пользователь с указанным id не найден, вернется ошибка NotFound.

	// - Если почта администратора не найдена — NotFound.
	// - Если пароль администратора не совпал — Unauthenticated.

	// При любых других ошибках – Internal.
	// Показатель успеха — информация о найденном пользователе.

	s.mu.RLock()
	_, err := utils.ComparePassword(ctx, req.Email, req.Password, true)
	s.mu.RUnlock()
	if err != nil {
		return nil, err
	}

	s.mu.RLock()
	res, err := db.GetUserById(ctx, req.Id)
	s.mu.RUnlock()
	if err != nil {
		return nil, utils.CompareErrAndErrNotFound(err)
	}

	return utils.ParseToResult(res), nil
}

// Обновляет информацию о пользователе по заданным полям.
func (s *AdminService) UpdateUser(ctx context.Context, req *pb.UpdateRequestAdmin) (*pb.Empty, error) {
	// ADMIN могут изменять только USER.
	// MAIN_ADMIN может изменять ADMIN и USER.
	// Гарантируется, что MAIN_ADMIN будет единственным.

	// - Если ADMIN попытается изменить ADMIN или MAIN_ADMIN, то вернется ошибка PermissionDenied.
	// - Если изменяемый пользователь не найден — NotFound.
	// - Если администратор попробует изменить самого себя — InvalidArgument.

	// - Если почта администратора не найдена — NotFound.
	// - Если пароль администратора не совпал — Unauthenticated.

	// При любых других ошибках – Internal.
	// Показатель успеха — отсутствие ошибки.

	s.mu.RLock()
	admin, err := utils.ComparePassword(ctx, req.AdminEmail, req.AdminPassword, true)
	s.mu.RUnlock()
	if err != nil {
		return nil, err
	}

	if (req.Role == pb.Role_MAIN_ADMIN || req.Role == pb.Role_ADMIN) && admin.Role == pb.Role_ADMIN.String() {
		return nil, status.Error(codes.PermissionDenied, "администратор может изменять только обычных пользователей")
	}

	if req.Email == admin.Email {
		return nil, status.Error(codes.InvalidArgument, "администратор не может изменять самого себя")
	}

	user := db.User{
		ID:         req.Id,
		Email:      req.Email,
		FirstName:  strings.ToLower(req.FirstName),
		SecondName: strings.ToLower(req.SecondName),
		Patronymic: strings.ToLower(*req.Patronymic),
	}

	s.mu.Lock()
	err = db.UpdateUser(ctx, user)
	s.mu.Unlock()
	if err != nil {
		return nil, utils.CompareErrAndErrNotFound(err)
	}

	return nil, nil
}

// Обновляет роль указанного пользователя.
func (s *AdminService) UpdateRole(ctx context.Context, req *pb.UpdateRoleRequestAdmin) (*pb.Empty, error) {
	// ADMIN не может изменять роли.
	// MAIN_ADMIN может изменять роли ADMIN и USER.
	// Т.к. гарантируется, что MAIN_ADMIN будет единственным, то при назначении кого-то на
	// MAIN_ADMIN, то первоначальный главный администратор станет обычным ADMIN.

	// - Если ADMIN попытается изменить чью-то роль, то вернется ошибка PermissionDenied.
	// - Если администратор попробует изменить самого себя — InvalidArgument.
	// - Если изменяемый пользователь не найден — NotFound.

	// - Если почта администратора не найдена — NotFound.
	// - Если пароль администратора не совпал — Unauthenticated.

	// При любых других ошибках – Internal.
	// Показатель успеха — отсутствие ошибки.

	s.mu.RLock()
	admin, err := utils.ComparePassword(ctx, req.AdminEmail, req.AdminPassword, true)
	s.mu.RUnlock()
	if err != nil {
		return nil, err
	}

	if admin.Role == pb.Role_ADMIN.String() {
		return nil, status.Error(codes.PermissionDenied, "администратор не может изменять роль пользователей")
	}

	s.mu.RLock()
	user, err := db.GetUserById(ctx, req.Id)
	s.mu.RUnlock()
	if err != nil {
		return nil, utils.CompareErrAndErrNotFound(err)
	}

	if user.Email == admin.Email {
		return nil, status.Error(codes.InvalidArgument, "администратор не может изменять собственную роль")
	}

	user.Role = req.NewRole.String()

	s.mu.Lock()
	err = db.UpdateUser(ctx, user)
	s.mu.Unlock()
	if err != nil {
		if errors.Is(err, db.ErrBadRequest) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, utils.CompareErrAndErrNotFound(err)
	}

	if req.NewRole == pb.Role_MAIN_ADMIN {
		admin.Role = pb.Role_ADMIN.String()

		s.mu.Lock()
		err = db.UpdateUser(ctx, admin)
		s.mu.Unlock()
		if err != nil {
			return nil, utils.CompareErrAndErrNotFound(err)
		}
	}
	return nil, nil
}

// Удаляет указанного пользователя.
func (s *AdminService) DeleteUser(ctx context.Context, req *pb.DeleteRequestAdmin) (*pb.Empty, error) {
	// ADMIN могут удалять только USER.
	// MAIN_ADMIN может удалять ADMIN и USER.
	// Гарантируется, что MAIN_ADMIN будет единственным.

	// - Если ADMIN попытается удалить ADMIN или MAIN_ADMIN, то вернется ошибка PermissionDenied.
	// - Если администратор попробует удалить самого себя — InvalidArgument.
	// - Если ни один пользователь с указанной почтой не найден, вернется ошибка NotFound.

	// - Если почта администратора не найдена — NotFound.
	// - Если пароль администратора не совпал — Unauthenticated.

	// При любых других ошибках – Internal.
	// Показатель успеха — информация о всех найденных пользователях.

	s.mu.RLock()
	admin, err := utils.ComparePassword(ctx, req.AdminEmail, req.AdminPassword, true)
	s.mu.RUnlock()
	if err != nil {
		return nil, err
	}

	s.mu.RLock()
	user, err := db.GetUserById(ctx, req.Id)
	s.mu.RUnlock()
	if err != nil {
		return nil, utils.CompareErrAndErrNotFound(err)
	}

	if user.Email == admin.Email {
		return nil, status.Error(codes.InvalidArgument, "администратор не может удалять самого себя")
	}

	if (user.Role == pb.Role_MAIN_ADMIN.String() || user.Role == pb.Role_ADMIN.String()) && admin.Role == pb.Role_ADMIN.String() {
		return nil, status.Error(codes.PermissionDenied, "администратор может удалять только обычных пользователей")
	}

	s.mu.Lock()
	err = db.DeleteUser(ctx, req.Id)
	s.mu.Unlock()
	if err != nil {
		return nil, utils.CompareErrAndErrNotFound(err)
	}
	return nil, nil
}
