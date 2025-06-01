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

	// При любых других ошибках – Internal.
	// Показатель успеха — отсутствие ошибки.

	s.mu.RLock()
	_, err := utils.ComparePassword(req.AdminEmail, req.AdminPassword, true)
	s.mu.RUnlock()

	if err != nil {
		return nil, err
	}

	user := db.User{
		Email:      req.GetEmail(),
		FirstName:  req.GetFirstName(),
		SecondName: req.GetSecondName(),
		Patronymic: req.GetPatronymic(),
		Password:   req.GetPassword(),
	}

	s.mu.Lock()
	_, err = db.CreateUser(user)
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
func (s *AdminService) GetUser(ctx context.Context, req *pb.GetRequestAdmin) (*pb.GetResponseAdmin, error) {
	// - Если ни один пользователь с указанной почтой не найден, вернется ошибка NotFound.

	// - Если почта администратора не найдена — NotFound.
	// - Если пароль администратора не совпал — Unauthenticated.

	// При любых других ошибках – Internal.
	// Показатель успеха — информация о всех найденных пользователях.

	s.mu.RLock()
	_, err := utils.ComparePassword(req.AdminEmail, req.AdminPassword, true)
	s.mu.RUnlock()
	if err != nil {
		return nil, err
	}

	s.mu.RLock()
	res, err := db.GetUserByKey(&req.SortBy, req.SortKey)
	s.mu.RUnlock()
	if err != nil {
		return nil, utils.CompareErrAndErrNotFound(err)
	}

	var users []*pb.GetResponse
	for r := range res {
		users = append(users, utils.ParseToResult(res[r]))
	}

	return &pb.GetResponseAdmin{Users: users}, nil
}

// Обновляет информацию о пользователе по заданным полям.
func (s *AdminService) UpdateUser(ctx context.Context, req *pb.UpdateRequestAdmin) (*pb.Empty, error) {
	// ADMIN могут изменять только USER.
	// MAIN_ADMIN может изменять ADMIN и USER.
	// Гарантируется, что MAIN_ADMIN будет единственным.

	// - Если ADMIN попытается изменить ADMIN или MAIN_ADMIN, то вернется ошибка PermissionDenied.
	// - Если MAIN_ADMIN попытается изменить другого MAIN_ADMIN, то вернется Internal, т.к. должен
	// быть только один MAIN_ADMIN.
	// - Если изменяемый пользователь не найден — NotFound.
	// - Если новая почта изменяемого пользователя уже существует — BadRequest.

	// - Если почта администратора не найдена — NotFound.
	// - Если пароль администратора не совпал — Unauthenticated.

	// При любых других ошибках – Internal.
	// Показатель успеха — отсутствие ошибки.

	s.mu.RLock()
	admin, err := utils.ComparePassword(req.AdminEmail, req.AdminPassword, true)
	s.mu.RUnlock()
	if err != nil {
		return nil, err
	}

	if (req.Role == pb.Role_MAIN_ADMIN || req.Role == pb.Role_ADMIN) && admin.Role == pb.Role_ADMIN.String() {
		return nil, status.Error(codes.PermissionDenied, "Администратор может изменять только обычных пользователей")
	}

	if req.Role == pb.Role_MAIN_ADMIN && admin.Role == pb.Role_MAIN_ADMIN.String() {
		return nil, status.Error(codes.Internal, "Должен быть только один MAIN_ADMIN!!!")
	}

	user := db.User{
		Email:      req.GetEmail(),
		FirstName:  strings.ToLower(req.GetFirstName()),
		SecondName: strings.ToLower(req.GetSecondName()),
		Patronymic: strings.ToLower(req.GetPatronymic()),
	}

	s.mu.Lock()
	err = db.UpdateUser(req.GetId(), user)
	s.mu.Unlock()
	if err != nil {
		if errors.Is(err, db.ErrBadRequest) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, utils.CompareErrAndErrNotFound(err)
	}

	return nil, nil
}

// Удаляет указанного пользователя. Притом обычные администраторы могут удалять только
// пользователей, когда как администраторов может удалять только единственный главный администратор.
func (s *AdminService) DeleteUser(ctx context.Context, req *pb.DeleteRequestAdmin) (*pb.Empty, error) {
	// ADMIN могут удалять только USER.
	// MAIN_ADMIN может удалять ADMIN и USER.
	// Гарантируется, что MAIN_ADMIN будет единственным.

	// - Если ADMIN попытается удалить ADMIN или MAIN_ADMIN, то вернется ошибка PermissionDenied.
	// - Если MAIN_ADMIN попытается изменить другого MAIN_ADMIN, то вернется Internal, т.к. должен
	// быть только один MAIN_ADMIN.
	// - Если ни один пользователь с указанной почтой не найден, вернется ошибка NotFound.

	// - Если почта администратора не найдена — NotFound.
	// - Если пароль администратора не совпал — Unauthenticated.

	// При любых других ошибках – Internal.
	// Показатель успеха — информация о всех найденных пользователях.

	s.mu.RLock()
	admin, err := utils.ComparePassword(req.AdminEmail, req.AdminPassword, true)
	s.mu.RUnlock()
	if err != nil {
		return nil, err
	}

	s.mu.RLock()
	user, err := db.GetUserById(req.Id)
	s.mu.RUnlock()
	if err != nil {
		return nil, utils.CompareErrAndErrNotFound(err)
	}

	if (user.Role == pb.Role_MAIN_ADMIN.String() || user.Role == pb.Role_ADMIN.String()) && admin.Role == pb.Role_ADMIN.String() {
		return nil, status.Error(codes.PermissionDenied, "Администратор может удалять только обычных пользователей")
	}

	if user.Role == pb.Role_MAIN_ADMIN.String() && admin.Role == pb.Role_MAIN_ADMIN.String() {
		return nil, status.Error(codes.Internal, "Должен быть только один MAIN_ADMIN!!!")
	}

	s.mu.Lock()
	err = db.DeleteUser(req.Id)
	s.mu.Unlock()
	if err != nil {
		return nil, utils.CompareErrAndErrNotFound(err)
	}
	return nil, nil
}
