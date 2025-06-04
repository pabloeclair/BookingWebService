package centraluniversity.app.booking.services;

import java.util.List;

import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;

import centraluniversity.app.booking.models.exception.HttpStatusException;
import centraluniversity.app.booking.models.user.*;
import centraluniversity.app.booking.pb.*;
import centraluniversity.app.booking.repositories.BookingRepository;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import io.grpc.Status;
import io.grpc.StatusRuntimeException;
import jakarta.annotation.PostConstruct;
import jakarta.annotation.PreDestroy;

@Service
public class AdminUserService {
    
    private ManagedChannel channel;
    private AdminServiceGrpc.AdminServiceBlockingStub stub;
    private final BookingRepository bookingRepository;

    AdminUserService(BookingRepository bookingRepository) {
        this.bookingRepository = bookingRepository;
    }

    // TODO: unit-test
    /**
     * Преобразовывается protocol buffers GetResponse в GetUserDto.
     * @param user - GetResponse
     * @return user - GetUserDto
     */
    private GetUserDto parseToDto(GetResponse user) {
        return new GetUserDto(
            user.getId(),
            user.getEmail(),
            user.getFirstName(),
            user.getSecondName(),
            user.getPatronymic(),
            user.getPassword(),
            user.getRole()
        );
    }

    /* Подключение к gRPC серверу авторизации. */
    @PostConstruct
    public void connectToServer() {
        this.channel = ManagedChannelBuilder.forAddress("admin-service", 7002)
            .usePlaintext()
            .build();
        this.stub = AdminServiceGrpc.newBlockingStub(channel);
    }

    /* Отсоединение от gRPC сервера авторизации. */
    @PreDestroy
    public void shutdown() {
        if (this.channel != null) {
            this.channel.shutdown();
        }
    }

    /**
     * Регистрация нового пользователя администратором.
     * @param user - полная информация о пользователе
     * @throws HttpStatusException CONFLICT (почта уже существует), NOT_FOUND (почта админа не найдена) UNAUTHORIZED (пароль админа не совпадает), FORBIDDEN (отказано в доступе)
     */
    public void createUser(CreateUserDto user) throws HttpStatusException {

        CreateRequestAdmin req;
        if (user.getPatronymic() == null) {
            req = CreateRequestAdmin.newBuilder()
                .setEmail(user.getEmail())
                .setFirstName(user.getFirstName())
                .setSecondName(user.getSecondName())
                .setPassword(user.getPassword())
                .setAdminEmail(user.getAdminEmail())
                .setAdminPassword(user.getAdminPassword())
                .build();
        } else {
            req = CreateRequestAdmin.newBuilder()
                .setEmail(user.getEmail())
                .setFirstName(user.getFirstName())
                .setSecondName(user.getSecondName())
                .setPatronymic(user.getPatronymic())
                .setPassword(user.getPassword())
                .setAdminEmail(user.getAdminEmail())
                .setAdminPassword(user.getAdminPassword())
                .build();
        }

        try {
            this.stub.createUser(req);
        } catch (StatusRuntimeException e) {
            Status status = e.getStatus();
            switch (status.getCode()) {
                case NOT_FOUND:
                    throw new HttpStatusException(HttpStatus.NOT_FOUND, e.getMessage());
                case ALREADY_EXISTS:
                    throw new HttpStatusException(HttpStatus.CONFLICT, e.getMessage());
                case UNAUTHENTICATED:
                    throw new HttpStatusException(HttpStatus.UNAUTHORIZED, e.getMessage());
                case PERMISSION_DENIED:
                    throw new HttpStatusException(HttpStatus.FORBIDDEN, e.getMessage());
                default:
                    throw new HttpStatusException(HttpStatus.INTERNAL_SERVER_ERROR, e.getMessage());
            }
        } 
    }

    /**
     * Сортировка всех пользователей по ключу администратором.
     * @param email
     * @param key
     * @param sortBy
     * @param sortKey
     * @return пустой список или список с полной информацией о найденных пользователях
     * @throws HttpStatusException NOT_FOUND (почта админа не найдена), UNAUTHORIZED (пароль админа не совпадает), FORBIDDEN (отказано в доступе)
     */
    public GetUserAdminDto sortUser(String email, String key, By sortBy, String sortKey) throws HttpStatusException {

        GetRequestAdmin req = GetRequestAdmin.newBuilder()
            .setAdminEmail(email)
            .setAdminPassword(key)
            .setSortBy(sortBy)
            .setSortKey(sortKey)
            .build();

        GetResponseArray res;
        try {
            res = this.stub.sortUser(req);
        } catch (StatusRuntimeException e) {
            Status status = e.getStatus();
            switch (status.getCode()) {
                case NOT_FOUND:
                    throw new HttpStatusException(HttpStatus.NOT_FOUND, e.getMessage());
                case UNAUTHENTICATED:
                    throw new HttpStatusException(HttpStatus.UNAUTHORIZED, e.getMessage());
                case PERMISSION_DENIED:
                    throw new HttpStatusException(HttpStatus.FORBIDDEN, e.getMessage());
                default:
                    throw new HttpStatusException(HttpStatus.INTERNAL_SERVER_ERROR, e.getMessage());
            }
        }

        List<GetResponse> users = res.getUsersList();
        GetUserDto[] result = new GetUserDto[res.getUsersCount()];
        for (int i = 0; i < res.getUsersCount(); i++) {
            GetResponse resUser = users.get(i);
            result[i] = parseToDto(resUser);
        }

        return new GetUserAdminDto(result);
    }

    /**
     * Получение пользователя по его id администратором.
     * @param id
     * @param email
     * @param key
     * @return
     * @throws HttpStatusException NOT_FOUND (почта админа не найдена), UNAUTHORIZED (пароль админа не совпадает), FORBIDDEN (отказано в доступе)
     */
    public GetUserDto getUserById(int id, String email, String key) throws HttpStatusException {

        Email req = Email.newBuilder().setEmail(email).setPassword(key).setId(id).build();

        GetResponse res;
        try {
            res = this.stub.getUserById(req);
        } catch (StatusRuntimeException e) {
            Status status = e.getStatus();
            switch (status.getCode()) {
                case NOT_FOUND:
                    throw new HttpStatusException(HttpStatus.NOT_FOUND, e.getMessage());
                case UNAUTHENTICATED:
                    throw new HttpStatusException(HttpStatus.UNAUTHORIZED, e.getMessage());
                case PERMISSION_DENIED:
                    throw new HttpStatusException(HttpStatus.FORBIDDEN, e.getMessage());
                default:
                    throw new HttpStatusException(HttpStatus.INTERNAL_SERVER_ERROR, e.getMessage());
            }
        }

        return parseToDto(res); 
    }

    /**
     * Обновление общей информации о пользователе администратором.
     * @param user
     * @throws HttpStatusException CONFLICT (почта уже существует), NOT_FOUND (почта админа/user не найден(-а)), UNAUTHORIZED (пароль админа не совпадает), FORBIDDEN (отказано в доступе)
     */
    public void updateUser(UpdateUserDto user) throws HttpStatusException {

        UpdateRequestAdmin req;
        if (user.getPatronymic() == null) {
            req = UpdateRequestAdmin.newBuilder()
                .setEmail(user.getEmail())
                .setFirstName(user.getFirstName())
                .setSecondName(user.getSecondName())
                .setPassword(user.getPassword())
                .setAdminEmail(user.getAdminEmail())
                .setAdminPassword(user.getAdminPassword())
                .setId(user.getId())
                .build();
        } else {
            req = UpdateRequestAdmin.newBuilder()
                .setEmail(user.getEmail())
                .setFirstName(user.getFirstName())
                .setSecondName(user.getSecondName())
                .setPatronymic(user.getPatronymic())
                .setPassword(user.getPassword())
                .setAdminEmail(user.getAdminEmail())
                .setAdminPassword(user.getAdminPassword())
                .setId(user.getId())
                .build();
        }

        try {
            this.stub.updateUser(req);
        } catch (StatusRuntimeException e) {
            Status status = e.getStatus();
            switch (status.getCode()) {
                case NOT_FOUND:
                    throw new HttpStatusException(HttpStatus.NOT_FOUND, e.getMessage());
                case UNAUTHENTICATED:
                    throw new HttpStatusException(HttpStatus.UNAUTHORIZED, e.getMessage());
                case PERMISSION_DENIED:
                    throw new HttpStatusException(HttpStatus.FORBIDDEN, e.getMessage());
                case ALREADY_EXISTS:
                    throw new HttpStatusException(HttpStatus.CONFLICT, e.getMessage());
                default:
                    throw new HttpStatusException(HttpStatus.INTERNAL_SERVER_ERROR, e.getMessage());
            }
        }
    }

    /**
     * Изменение роли пользователя или администратора или передача прав.
     * @param id
     * @param user
     * @throws HttpStatusException BAD_REQUEST (администратор не может менять сам себя), NOT_FOUND (почта админа/user не найден(-а)) UNAUTHORIZED (пароль админа не совпадает), FORBIDDEN (отказано в доступе)
     */
    public void updateRole(int id, UpdateRoleUserDto user) throws HttpStatusException {

        UpdateRoleRequestAdmin req = UpdateRoleRequestAdmin.newBuilder()
            .setAdminEmail(user.getAdminEmail())
            .setAdminPassword(user.getAdminPassword())
            .setNewRole(Role.valueOf(user.getNewRole()))
            .setOldRole(Role.valueOf(user.getOldRole()))
            .setId(id).build();

        try {
            this.stub.updateRole(req);
        } catch (StatusRuntimeException e) {
            Status status = e.getStatus();
            switch (status.getCode()) {
                case NOT_FOUND:
                    throw new HttpStatusException(HttpStatus.NOT_FOUND, e.getMessage());
                case UNAUTHENTICATED:
                    throw new HttpStatusException(HttpStatus.UNAUTHORIZED, e.getMessage());
                case PERMISSION_DENIED:
                    throw new HttpStatusException(HttpStatus.FORBIDDEN, e.getMessage());
                case INVALID_ARGUMENT:
                    throw new HttpStatusException(HttpStatus.BAD_REQUEST, e.getMessage());
                default:
                    throw new HttpStatusException(HttpStatus.INTERNAL_SERVER_ERROR, e.getMessage());
            }
        }
    }

    /**
     * Удаление пользователя администратором.
     * @param id
     * @param email
     * @param key
     * @throws HttpStatusException NOT_FOUND (почта админа/user не найден(-а)), UNAUTHORIZED (пароль админа не совпадает), FORBIDDEN (отказано в доступе)
     */
    public void deleteUser(int id, String email, String key) throws HttpStatusException {

        DeleteRequestAdmin req = DeleteRequestAdmin.newBuilder()
            .setAdminEmail(email)
            .setAdminPassword(key)
            .setId(id)
            .build();
        
        try {
            this.stub.deleteUser(req);
        } catch (StatusRuntimeException e) {
            Status status = e.getStatus();
            switch (status.getCode()) {
                case NOT_FOUND:
                    throw new HttpStatusException(HttpStatus.NOT_FOUND, e.getMessage());
                case UNAUTHENTICATED:
                    throw new HttpStatusException(HttpStatus.UNAUTHORIZED, e.getMessage());
                case PERMISSION_DENIED:
                    throw new HttpStatusException(HttpStatus.FORBIDDEN, e.getMessage());
                default:
                    throw new HttpStatusException(HttpStatus.INTERNAL_SERVER_ERROR, e.getMessage());
            }
        }

        bookingRepository.deleteByUserId(id);
    }
}
