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

    public void createUser(CreateUserDto user) {

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
    }

    public GetUserAdminDto sortUser(String email, String key, By sortBy, String sortKey) {

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

    public GetUserDto getUserById(int id, String email, String key) {

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

    public void updateUser(UpdateUserDto user) {

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
                    throw new HttpStatusException(HttpStatus.BAD_REQUEST, e.getMessage());
                default:
                    throw new HttpStatusException(HttpStatus.INTERNAL_SERVER_ERROR, e.getMessage());
            }
        }
    }

    public void updateRole(int id, UpdateRoleUserDto user) {

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
                case ALREADY_EXISTS:
                    throw new HttpStatusException(HttpStatus.BAD_REQUEST, e.getMessage());
                default:
                    throw new HttpStatusException(HttpStatus.INTERNAL_SERVER_ERROR, e.getMessage());
            }
        }
    }

    public void deleteUser(int id, String email, String key) {

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
