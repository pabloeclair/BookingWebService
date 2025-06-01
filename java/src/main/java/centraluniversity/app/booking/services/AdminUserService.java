package centraluniversity.app.booking.services;

import java.util.List;

import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;

import centraluniversity.app.booking.models.exception.HttpStatusException;
import centraluniversity.app.booking.models.user.*;
import centraluniversity.app.booking.pb.AdminServiceGrpc;
import centraluniversity.app.booking.pb.*;
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

    private UserResponseDto parseToDto(UserResponse user) {
        return new UserResponseDto(
            user.getId(),
            user.getEmail(),
            user.getFirstName(),
            user.getSecondName(),
            user.getPatronymic(),
            user.getPassword(),
            user.getRole()
        );
    }

    @PostConstruct
    public void connectToServer() {
        this.channel = ManagedChannelBuilder.forAddress("admin", 7002)
            .usePlaintext()
            .build();
        this.stub = AdminServiceGrpc.newBlockingStub(channel);
    }

    @PreDestroy
    public void shutdown() {
        if (this.channel != null) {
            this.channel.shutdown();
        }
    }

    public UserResponseDto createUser(CreateUserDto user) {

        CreateRequestAdmin req;
        if (user.getPatronymic() == null || user.getPatronymic().isEmpty()) {
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

        UserResponse res;
        try {
            res = this.stub.createUser(req);
        } catch (StatusRuntimeException e) {
            Status status = e.getStatus();
            switch (status.getCode()) {
                case ALREADY_EXISTS:
                    throw new HttpStatusException(HttpStatus.BAD_REQUEST, e.getMessage());
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

    public GetUserDto getUser(String email, String key, By sortBy, String sortKey) {

        GetRequestAdmin req = GetRequestAdmin.newBuilder()
            .setAdminEmail(email)
            .setAdminPassword(key)
            .setSortBy(sortBy)
            .setSortKey(sortKey)
            .build();

        GetUserResponseAdmin res;
        try {
            res = this.stub.getUser(req);
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

        List<UserResponse> users = res.getUsersList();
        UserResponseDto[] result = new UserResponseDto[res.getUsersCount()];
        for (int i = 0; i < res.getUsersCount(); i++) {
            UserResponse resUser = users.get(i);
            result[i] = parseToDto(resUser);
        }

        return new GetUserDto(result);
    }

    public UserResponseDto updateUser(Integer id, UpdateUserDto user) {

        UpdateRequestAdmin req;
        if (user.getPatronymic() == null || user.getPatronymic().isEmpty()) {
            req = UpdateRequestAdmin.newBuilder()
                .setEmail(user.getEmail())
                .setFirstName(user.getFirstName())
                .setSecondName(user.getSecondName())
                .setPassword(user.getPassword())
                .setAdminEmail(user.getAdminEmail())
                .setAdminPassword(user.getAdminPassword())
                .setId(id)
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
                .setId(id)
                .build();
        }

        UserResponse res;
        try {
            res = this.stub.updateUser(req);
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
                    throw new HttpStatusException(HttpStatus.BAD_REQUEST, e.getMessage());
                default:
                    throw new HttpStatusException(HttpStatus.INTERNAL_SERVER_ERROR, e.getMessage());
            }
        }

        return parseToDto(res);
    }

    public void deleteUser(Integer id, String email, String key) {

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
    }
}
