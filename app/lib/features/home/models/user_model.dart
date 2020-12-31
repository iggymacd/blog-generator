// import 'package:cloud_firestore/cloud_firestore.dart';
import 'package:freezed_annotation/freezed_annotation.dart';

import '../../../../core/types/gender.dart';
import '../../../../core/types/vital_status.dart';
import '../../../../core/util/mapper.dart';
import '../../domain/entities/character.dart';

part 'user_model.freezed.dart';
part 'character_model.g.dart';

class UserModel {
// class Model {
  // static const ID = "id";
  // static const DESCRIPTION = "description";
  // static const CART = "cart";
  // static const USER_ID = "userId";
  // static const TOTAL = "total";
  // static const STATUS = "status";
  // static const CREATED_AT = "createdAt";
  // static const UPDATED_AT = "updatedAt";
 
	static const CREATED_AT = "createdAt"; 
	static const CREATED_BY = "createdBy"; 
	static const GENDER = "gender"; 
	static const UPDATED_AT = "updatedAt"; 
	static const UPDATED_BY = "updatedBy"; 

  // String _id;
  // String _description;
  // String _userId;
  // String _status;
    // range .Fields
    //     String .Name;
    // end
   
	String _createdAt; 
	String _createdBy; 
	String _gender; 
	String _updatedAt; 
	String _updatedBy; 

  int _createdAt;
  int _updatedAt;
  // int _total;

//  getters
  // String get id => _id;

  // String get description => _description;

  // String get userId => _userId;

  // String get status => _status;

  // int get total => _total;

  int get createdAt => _createdAt;

  int get updatedAt => _updatedAt;

  // public variable
  // List cart;

  // OrderModel.fromSnapshot(DocumentSnapshot snapshot) {
  //   _id = snapshot.data[ID];
  //   _description = snapshot.data[DESCRIPTION];
  //   _total = snapshot.data[TOTAL];
  //   _status = snapshot.data[STATUS];
  //   _userId = snapshot.data[USER_ID];
  //   _createdAt = snapshot.data[CREATED_AT];
  //   cart = snapshot.data[CART];
  // }
}
