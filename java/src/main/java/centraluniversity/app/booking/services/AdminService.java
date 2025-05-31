package centraluniversity.app.booking.services;

import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;

import centraluniversity.app.booking.models.UserResponseDto;
import centraluniversity.app.booking.models.admin.CreateUserDto;
import centraluniversity.app.booking.models.exception.HttpStatusException;
import centraluniversity.app.booking.pb.AdminServiceGrpc;
import centraluniversity.app.booking.pb.CreateRequestAdmin;
import centraluniversity.app.booking.pb.UserResponse;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import io.grpc.Status;
import io.grpc.StatusRuntimeException;
import jakarta.annotation.PostConstruct;
import jakarta.annotation.PreDestroy;

@Service
public class AdminService {
    
    private ManagedChannel channel;
    private AdminServiceGrpc.AdminServiceBlockingStub stub;

    @PostConstruct
    public void connectToServer() {
        this.channel = ManagedChannelBuilder.forAddress("admin", 7001)
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
            if (status.getCode() == Status.Code.ALREADY_EXISTS) {
                throw new HttpStatusException(HttpStatus.BAD_REQUEST, e.getMessage());
            } else if (status.getCode() == Status.Code.UNAUTHENTICATED) {
                throw new HttpStatusException(HttpStatus.UNAUTHORIZED, e.getMessage());
            } else if (status.getCode() == Status.Code.PERMISSION_DENIED) {
                throw new HttpStatusException(HttpStatus.FORBIDDEN, e.getMessage());
            }
            throw new HttpStatusException(HttpStatus.INTERNAL_SERVER_ERROR, e.getMessage());
        } 

        return new UserResponseDto(
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
