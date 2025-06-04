package centraluniversity.app.booking.services;

import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;

import centraluniversity.app.booking.models.exception.HttpStatusException;
import centraluniversity.app.booking.models.user.SignupUserDto;
import centraluniversity.app.booking.models.user.GetUserDto;
import centraluniversity.app.booking.models.user.IdDto;
import centraluniversity.app.booking.pb.AuthenticationGrpc;
import centraluniversity.app.booking.pb.Email;
import centraluniversity.app.booking.pb.GetResponse;
import centraluniversity.app.booking.pb.Id;
import centraluniversity.app.booking.pb.SignupRequest;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import io.grpc.Status;
import io.grpc.StatusRuntimeException;
import jakarta.annotation.PostConstruct;
import jakarta.annotation.PreDestroy;

@Service
public class AuthService {
    
    private ManagedChannel channel;
    private AuthenticationGrpc.AuthenticationBlockingStub stub;

    /* Подключение к gRPC серверу авторизации. */
    @PostConstruct
    public void connectToServer() {
        this.channel = ManagedChannelBuilder.forAddress("auth-service", 7001)
            .usePlaintext()
            .build();
        this.stub = AuthenticationGrpc.newBlockingStub(channel);
    }

    /* Отсоединение от gRPC сервера авторизации. */
    @PreDestroy
    public void shutdown() {
        if (this.channel != null) {
            this.channel.shutdown();
        }
    }

    /**
     * Регистрация нового пользователя.
     * @param user - полная информация о пользователе
     * @return id пользователя
     * @throws HttpStatusException CONFLICT (почта уже существует)
     */
    public IdDto createUser(SignupUserDto user) throws HttpStatusException {

        SignupRequest req;
        if (user.getPatronymic() == null || user.getPatronymic().isEmpty()) {
            req = SignupRequest.newBuilder()
                .setEmail(user.getEmail())
                .setFirstName(user.getFirstName())
                .setSecondName(user.getSecondName())
                .setPassword(user.getPassword())
                .build();
        } else {
            req = SignupRequest.newBuilder()
                .setEmail(user.getEmail())
                .setFirstName(user.getFirstName())
                .setSecondName(user.getSecondName())
                .setPatronymic(user.getPatronymic())
                .setPassword(user.getPassword())
                .build();
        }

        Id res;
        try {
            res = this.stub.signupUser(req);
        } catch (StatusRuntimeException e) {
            Status status = e.getStatus();
            if (status.getCode() == Status.Code.ALREADY_EXISTS) {
                throw new HttpStatusException(HttpStatus.CONFLICT, e.getMessage());
            } 
            throw new HttpStatusException(HttpStatus.INTERNAL_SERVER_ERROR, e.getMessage());
        } 
        
        return new IdDto(res.getId());  
    }

    /**
     * Авторизация и получение информации о пользователе.
     * @param email
     * @return полная информация о пользователе
     * @throws HttpStatusException NOT_FOUND (почта не найдена), UNAUTHORIZED (пароль не совпадает), FORBIDDEN (доступ запрещен)
     */
    public GetUserDto getUserByEmail(String email, String password) throws HttpStatusException {

        Email req = Email.newBuilder().setEmail(email).setPassword(password).build();
        GetResponse res;
        try {
            res = this.stub.loginUser(req);
        } catch (StatusRuntimeException e) {
            Status status = e.getStatus();
            switch (status.getCode()) {
                case NOT_FOUND:
                    throw new HttpStatusException(HttpStatus.NOT_FOUND, e.getMessage());
                case UNAUTHENTICATED:
                    throw new HttpStatusException(HttpStatus.UNAUTHORIZED, e.getMessage());
                default:
                    throw new HttpStatusException(HttpStatus.INTERNAL_SERVER_ERROR, e.getMessage());
            }
        }
        
        return new GetUserDto(
            res.getId(),
            res.getEmail(),
            res.getFirstName(),
            res.getSecondName(),
            res.getPatronymic(),
            res.getPassword(),
            res.getRole()
        ); 
    }

}
